package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server-database/types"
	"strconv"
)

func handleImages(w http.ResponseWriter, r *http.Request) {
	images, _ := fetchImages()
	indent := getIndentParam(r)
	b, err := json.MarshalIndent(images, "", indent)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getIndentParam(r *http.Request) string {
	params := r.URL.Query()

	indentSize, err := strconv.Atoi(params.Get("indent"))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	indent := ""
	for i := 0; i < indentSize; i++ {
		indent += " "
	}
	return indent
}

func fetchImages() ([]types.Image, error) {
	images := []types.Image{
		{
			Title:   "Sunset",
			AltText: "Clouds at sunset",
			URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
		{
			Title:   "Mountain",
			AltText: "A mountain at sunset",
			URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
	}

	return images, nil
}

func main() {

	http.HandleFunc("/images.json", handleImages)

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}
