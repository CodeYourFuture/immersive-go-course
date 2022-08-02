package main

import (
	"encoding/json"
	"net/http"
)

type Image struct {
	Title   string
	AltText string
	Url     string
}

func main() {
	data := []Image{
		{"Sunset", "Clouds at sunset", "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80"},
	}

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		// Indicate that what follows will be JSON
		w.Header().Add("Content-Type", "text/json")
		// Convert our data file to a byte-array for writing back in a response
		b, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		// Send it back!
		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)
}
