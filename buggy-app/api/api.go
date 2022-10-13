package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
)

type ApiService struct {
	util.Service
}

type Config struct {
	Port int
	Log  *log.Logger
}

func (as *ApiService) handleMyNotes(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (as *ApiService) Run(ctx context.Context, config Config) error {
	listen := fmt.Sprintf("localhost:%d", config.Port)

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

func NewApiService() ApiService {
	return ApiService{}
}
