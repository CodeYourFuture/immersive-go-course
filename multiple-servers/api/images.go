package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func FetchImages(conn *pgx.Conn) ([]Image, error) {
	// Send a query to the database, returning raw rows
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	// Handle query errors
	if err != nil {
		return nil, fmt.Errorf("unable to query database: [%w]", err)
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

func AddImage(conn *pgx.Conn, r *http.Request) (*Image, error) {
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
