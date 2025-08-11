package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Config struct {
	DatabaseURL string
	Port        int
}

func Run(config Config) error {
	conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	// Defer closing the connection to when main function exits
	defer conn.Close(context.Background())

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.EscapedPath())

		// Grab the indent query param early
		indent := r.URL.Query().Get("indent")

		var response []byte
		var responseErr error
		if r.Method == "POST" {
			// Add new image to the database
			image, err := AddImage(conn, r)
			if err != nil {
				log.Println(err.Error())
				// We don't expose our internal errors (i.e. the contents of err) directly to the user for a few reasons:
				//  1. It may leak private information (e.g. a database connection string, which may even include a password!), which may be a security risk.
				//  2. It probably isn't useful to them to know.
				//  3. It may contain confusing terminology which may be embarrassing or confusing to expose.
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			response, responseErr = MarshalWithIndent(image, indent)
		} else {
			// Fetch images from the database
			images, err := FetchImages(conn)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			response, responseErr = MarshalWithIndent(images, indent)
		}

		if responseErr != nil {
			log.Println(responseErr.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		// Indicate that what follows will be JSON
		w.Header().Add("Content-Type", "text/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		// Send it back!
		w.Write(response)
	})

	log.Printf("port: %d\n", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
