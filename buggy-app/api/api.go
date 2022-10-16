package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
)

type ApiService struct {
	util.Service

	authClient auth.Client
}

type Config struct {
	Port           int
	Log            *log.Logger
	AuthServiceUrl string
}

func (as *ApiService) handleMyNotes(w http.ResponseWriter, r *http.Request) {
	id, passwd, ok := r.BasicAuth()
	// Malformed basic auth is not OK
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Use the auth client to check if this id/password combo is approved
	result, err := as.authClient.Verify(id, passwd)
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

	// All good!
	w.Header().Add("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}

func (as *ApiService) Run(ctx context.Context, config Config) error {
	listen := fmt.Sprintf("localhost:%d", config.Port)

	client, err := auth.NewClient(ctx, config.AuthServiceUrl)
	if err != nil {
		return err
	}
	as.authClient = client

	mux := new(http.ServeMux)
	mux.HandleFunc("/1/my/notes.json", as.handleMyNotes)

	server := &http.Server{Addr: listen, Handler: mux}

	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = server.ListenAndServe()
	}()

	config.Log.Printf("api service: listening: %s", listen)

	<-ctx.Done()
	server.Shutdown(context.TODO())

	wg.Wait()
	return runErr
}

func NewApiService() *ApiService {
	return &ApiService{}
}
