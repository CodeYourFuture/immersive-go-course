package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"golang.org/x/net/context"
)

func main() {
	port := flag.Int("port", 8080, "port the server will listen on")
	flag.Parse()

	// Get the postgres password from a file supplied in an environment variable
	// TODO: it would be better for this to come from DATABASE_URL or to "figure out"
	// the best auth params from environment variables
	passwd, err := util.ReadPasswdFile()
	if err != nil {
		log.Fatal(err)
	}

	// The NotifyContext will signal Done when these signals are sent, allowing the server
	// to shutdown gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	as := auth.NewAuthService()
	if err := as.Run(ctx, auth.Config{
		Port:        *port,
		DatabaseUrl: fmt.Sprintf("postgres://postgres:%s@postgres:5432/auth", passwd),
		Log:         log.Default(),
	}); err != nil {
		log.Fatal(err)
	}
}
