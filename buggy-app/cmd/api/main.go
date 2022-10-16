package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/api"
	"golang.org/x/net/context"
)

func main() {
	port := flag.Int("port", 8090, "port the server will listen on")
	flag.Parse()

	// The NotifyContext will signal Done when these signals are sent, allowing the server
	// to shutdown gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	as := api.NewApiService()
	if err := as.Run(ctx, api.Config{
		Port: *port,
		Log:  log.Default(),
	}); err != nil {
		log.Fatal(err)
	}
}
