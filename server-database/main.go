package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Image struct {
	Title   string
	AltText string
	Url     string
}

func main() {
	data := []Image{
		{"Sunset", "Clouds at sunset", "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"},
		{"Mountain", "A mountain at sunset", "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"},
	}

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		indent := r.URL.Query().Get("indent")
		// Convert our data file to a byte-array for writing back in a response
		var b []byte
		var marshalErr error
		// Allow up to 10 characters of indent
		if i, err := strconv.Atoi(indent); err == nil && i > 0 && i <= 10 {
			b, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", i))
		} else {
			b, marshalErr = json.Marshal(data)
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
