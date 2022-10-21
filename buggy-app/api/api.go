package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/api/model"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authuserctx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	httplogger "github.com/gleicon/go-httplogger"
)

// DbClient is for talking to the database
type DbClient interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Close()
}

type Config struct {
	Port           int
	Log            *log.Logger
	AuthServiceUrl string
	DatabaseUrl    string
}

type Service struct {
	config     Config
	authClient auth.Client
	pool       DbClient
}

func New(config Config) *Service {
	return &Service{
		config: config,
	}
}

// HTTP handler for getting notes for a particular user
func (as *Service) handleMyNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get the authenticated user from the context -- this will have been written earlier
	owner, ok := authuserctx.FromAuthenticatedContext(ctx)
	if !ok {
		as.config.Log.Printf("api: route handler reached with invalid auth context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

	// Use the "model" layer to get a list of the owner's notes
	notes, err := model.GetNotesForOwner(ctx, as.pool, owner)
	if err != nil {
		fmt.Printf("api: GetNotesForOwner failed: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	response := struct {
		Notes model.Notes `json:"notes"`
	}{
		Notes: notes,
	}

	// Convert the []Row into JSON
	res, err := util.MarshalWithIndent(response, "")
	if err != nil {
		fmt.Printf("api: response marshal failed: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Write it back out!
	w.Header().Add("Content-Type", "text/json")
	w.Write(res)
}

// HTTP handler for getting notes for a particular user
func (as *Service) handleMyNoteById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get the authenticated user from the context -- this will have been written earlier
	_, ok := authuserctx.FromAuthenticatedContext(ctx)
	if !ok {
		as.config.Log.Printf("api: route handler reached with invalid auth context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

	// The URL.Path will be something like /1/my/notes/abc123.json.
	// path.Base strips everything but "abc123.json". We then Replace out the ".json" to give us
	// just the ID.
	id := strings.Replace(path.Base(r.URL.Path), ".json", "", 1)
	if id == "" {
		fmt.Printf("api: no ID supplied: url path %v\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	// Use the "model" layer to get a list of the owner's notes
	note, err := model.GetNoteById(ctx, as.pool, id)
	if err != nil {
		fmt.Printf("api: GetNoteById failed: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	response := struct {
		Note model.Note `json:"note"`
	}{
		Note: note,
	}

	// Convert the []Row into JSON
	res, err := util.MarshalWithIndent(response, "")
	if err != nil {
		fmt.Printf("api: response marshal failed: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Write it back out!
	w.Header().Add("Content-Type", "text/json")
	w.Write(res)
}

// Set up routes -- this can be used in tests to set up simple HTTP handling
// rather than running the whole server.
func (as *Service) Handler() http.Handler {
	mux := new(http.ServeMux)
	mux.HandleFunc("/1/my/note/", as.wrapAuth(as.authClient, as.handleMyNoteById))
	mux.HandleFunc("/1/my/notes.json", as.wrapAuth(as.authClient, as.handleMyNotes))
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

	// Connect to the Auth service via the AuthClient
	client, err := auth.NewClient(ctx, as.config.AuthServiceUrl)
	if err != nil {
		return err
	}
	as.authClient = client

	// mux is the root Handler
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

	// Wait for a signal to shut down...
	<-ctx.Done()
	// ... and then do it as gracefully as possible.
	server.Shutdown(context.TODO())

	wg.Wait()
	return runErr
}
