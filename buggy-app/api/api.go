package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authuserctx"
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

type noteRow struct {
	Id       string    `json:"id"`
	Owner    string    `json:"owner"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
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

		ctx = authuserctx.NewAuthenticatedContext(ctx, id)
		handler(w, r.WithContext(ctx))
	}
}

// The special type interface{} allows us to take _any_ value, not just one of a specific type.
func marshalWithIndent(data interface{}, indent string) ([]byte, error) {
	// Convert images to a byte-array for writing back in a response
	var b []byte
	var marshalErr error
	// Allow up to 10 characters of indent
	if i, err := strconv.Atoi(indent); err == nil && i > 0 && i <= 10 {
		b, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", i))
	} else {
		b, marshalErr = json.Marshal(data)
	}
	if marshalErr != nil {
		return nil, fmt.Errorf("could not marshal data: [%w]", marshalErr)
	}
	return b, nil
}

func (as *Service) handleMyNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := authuserctx.FromAuthenticatedContext(ctx)
	if !ok {
		as.config.Log.Printf("api: route handler reached with invalid auth context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

	queryRows, err := as.pool.Query(ctx, "SELECT id, owner, content, created, modified FROM public.note")
	if err != nil {
		fmt.Printf("api: could not query notes: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	defer queryRows.Close()

	rows := []noteRow{}
	for queryRows.Next() {
		row := noteRow{}
		err = queryRows.Scan(&row.Id, &row.Owner, &row.Content, &row.Created, &row.Modified)
		if err != nil {
			fmt.Printf("api: query scan failed: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if row.Owner == userId {
			rows = append(rows, row)
		}
	}

	if queryRows.Err() != nil {
		fmt.Printf("api: query read failed: %v\n", queryRows.Err())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	res, err := marshalWithIndent(rows, "")
	if err != nil {
		fmt.Printf("api: response marshal failed: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
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
