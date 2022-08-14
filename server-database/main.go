package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

func (img Image) String() string {
	return fmt.Sprintf("%s (%s): %s", img.Title, img.AltText, img.URL)
}

func fetchImages(conn *pgx.Conn) ([]Image, error) {
	// Send a query to the database, returning raw rows
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
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
			return nil, fmt.Errorf("unable to read from database: %w", err)
		}
		// Append this as a new Image to the images slice
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}

	return images, nil
}

func addImage(conn *pgx.Conn, r *http.Request) (*Image, error) {
	// Read the request body into a bytes slice
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body: [%w]", err)
	}

	// Parse the body JSON into an image struct
	var image Image
	err = json.Unmarshal(body, &image)
	if err != nil {
		return nil, fmt.Errorf("could not parse request body: [%w]", err)
	}

	// Insert it into the database
	_, err = conn.Exec(context.Background(), "INSERT INTO public.images(title, url, alt_text) VALUES ($1, $2, $3)", image.Title, image.URL, image.AltText)
	if err != nil {
		return nil, fmt.Errorf("could not insert image: [%w]", err)
	}

	return &image, nil
}

// The special type interface{} allows us to take _any_ value, not just one of a specific type.
// This means we can re-use this function for both a single image _and_ a slice of multiple images.
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

	// Create the handler function that serves the images JSON
	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		// Grab the indent query param early
		indent := r.URL.Query().Get("indent")

		var response []byte
		var responseErr error
		if r.Method == "POST" {
			// Add new image to the database
			image, err := addImage(conn, r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				// We don't expose our internal errors (i.e. the contents of err) directly to the user for a few reasons:
				//  1. It may leak private information (e.g. a database connection string, which may even include a password!), which may be a security risk.
				//  2. It probably isn't useful to them to know.
				//  3. It may contain confusing terminology which may be embarrassing or confusing to expose.
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			response, responseErr = marshalWithIndent(image, indent)
		} else {
			// Fetch images from the database
			images, err := fetchImages(conn)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			response, responseErr = marshalWithIndent(images, indent)
		}

		if responseErr != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		// Indicate that what follows will be JSON
		w.Header().Add("Content-Type", "text/json")
		// Send it back!
		w.Write(response)
	})
	http.ListenAndServe(":8080", nil)
}
