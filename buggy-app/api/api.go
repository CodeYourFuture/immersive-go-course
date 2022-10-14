package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authctx"

	httplogger "github.com/gleicon/go-httplogger"
)

type Config struct {
	Port           int
	Log            *log.Logger
	AuthServiceUrl string
}

type Service struct {
	util.Service

	config     Config
	authClient auth.Client
}

func New(config Config) *Service {
	return &Service{
		config: config,
	}
}

func (as *Service) wrapAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		id, passwd, ok := r.BasicAuth()
		// Malformed basic auth is not OK
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Use the auth client to check if this id/password combo is approved
		result, err := as.authClient.Verify(ctx, id, passwd)
		if err != nil {
			log.Printf("api: verify error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Unless we get an Allow, say no
		if result.State != auth.StateAllow {
			log.Printf("api: verify denied: id %v\n", id)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = authctx.NewAuthenticatedContext(ctx, id)
		handler(w, r.WithContext(ctx))
	}
}

func (as *Service) handleMyNotes(w http.ResponseWriter, r *http.Request) {
	_, ok := authctx.FromAuthenticatedContext(r.Context())
	if !ok {
		as.config.Log.Printf("api: route handler reached with invalid auth context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	w.Header().Add("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}

// Set up routes -- this can be used in tests to set up simple HTTP handling
// rather than running the whole server.
func (as *Service) Handler() http.Handler {
	mux := new(http.ServeMux)
	mux.HandleFunc("/1/my/notes.json", as.wrapAuth(as.handleMyNotes))
	return httplogger.HTTPLogger(mux)
}

func (as *Service) Run(ctx context.Context) error {
	listen := fmt.Sprintf("localhost:%d", as.config.Port)

	client, err := auth.NewClient(ctx, as.config.AuthServiceUrl)
	if err != nil {
		return err
	}
	as.authClient = client

	mux := as.Handler()
	server := &http.Server{Addr: listen, Handler: mux}

	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = server.ListenAndServe()
	}()

	as.config.Log.Printf("api service: listening: %s", listen)

	<-ctx.Done()
	server.Shutdown(context.TODO())

	wg.Wait()
	return runErr
}
