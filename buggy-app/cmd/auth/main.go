package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"golang.org/x/net/context"
)

func main() {
	port := flag.Int("port", 8080, "port the server will listen on")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	as := auth.NewAuthService()
	if err := as.Run(ctx, auth.Config{
		Port: *port,
		Log:  *log.Default(),
	}); err != nil {
		log.Fatal(err)
	}
}
