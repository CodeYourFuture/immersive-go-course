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
	"github.com/jackc/pgx/v5/pgxpool"

	httplogger "github.com/gleicon/go-httplogger"
)

type Config struct {
	Port           int
	Log            *log.Logger
	AuthServiceUrl string
	DatabaseUrl    string
}

type Service struct {
	util.Service

	config     Config
	authClient auth.Client
	pool       *pgxpool.Pool
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
	listen := fmt.Sprintf(":%d", as.config.Port)

	// Connect to the database via a "pool" of connections, allowing concurrency
	pool, err := pgxpool.New(ctx, as.config.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	defer pool.Close()
	// Add the pool to the the service
	as.pool = pool

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
