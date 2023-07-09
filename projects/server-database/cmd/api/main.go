package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"time"

	"server-database/cmd/api/server"

	v1 "server-database/cmd/api/v1"

	_ "github.com/lib/pq"
)

func main() {
	svr := &server.Server{}
	if err := svr.MountLogger(); err != nil {
		os.Exit(1)
	}

	db, err := svr.MountDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	if err := svr.MountImageService(); err != nil {
		os.Exit(1)
	}
	mux := http.NewServeMux()
	v1.Register(mux, svr)

	mux.Handle("/metrics", promhttp.Handler())

	// TODO: move address to env file
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	svr.Logger.Printf("starting server on: %q", "8080")
	if err := httpServer.ListenAndServe(); err != nil {
		svr.Logger.Printf("error starting the server: %v", err)
		os.Exit(1)
	}
}
