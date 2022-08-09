package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string
	AltText string
	Url     string
}

func fetchImages(conn *pgx.Conn) ([]Image, error) {
	// Send a query to the database, returning raw rows
	rows, err := conn.Query(context.Background(), `"SELECT title, url, alt_text FROM public.images"`)
	// Handle query errors
	if err != nil {
		return []Image{}, fmt.Errorf("unable to query database: [%w]", err)
	}

	// Create slice to contain the images
	var images []Image
	// Iterate through each row to extract the data
	for rows.Next() {
		var title, url, altText string
		// Extract the data, passing pointers so the data can be updated in place
		err = rows.Scan(&title, &url, &altText)
		if err != nil {
			break
		}
		// Append this as a new Image to the images slice
		images = append(images, Image{Title: title, Url: url, AltText: altText})
	}

	return images, nil
}

func main() {
	// Check that DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		fmt.Fprintf(os.Stderr, "DATABASE_URL not set\n")
		os.Exit(1)
	}

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// Handle a possible connection error
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// Defer closing the connection to when main function exits
	defer conn.Close(context.Background())

	// Fetch images from the database
	images, err := fetchImages(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get images: %v\n", err)
		os.Exit(1)
	}

	// Create the handler function that serves the images JSON
	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		indent := r.URL.Query().Get("indent")
		// Convert images to a byte-array for writing back in a response
		var b []byte
		var marshalErr error
		// Allow up to 10 characters of indent
		if i, err := strconv.Atoi(indent); err == nil && i > 0 && i <= 10 {
			b, marshalErr = json.MarshalIndent(images, "", strings.Repeat(" ", i))
		} else {
			b, marshalErr = json.Marshal(images)
		}
		if marshalErr != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		// Indicate that what follows will be JSON
		w.Header().Add("Content-Type", "text/json")
		// Send it back!
		w.Write(b)
	})
	http.ListenAndServe(":8080", nil)
}
