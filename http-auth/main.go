package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Indicate that we are sending back HTML
		w.Header().Add("Content-Type", "text/html")
		// Write the doctype and opening tag regardless of method
		w.Write([]byte("<!DOCTYPE html><html>"))
		// If the request is POSTing data, return what they sent back
		if r.Method == "POST" {
			// The request (r) body is an io.Reader and the response (w) is a writer
			// so we can stream one directly into the other in chunks.
			// We ignore the output of io.Copy and just handle the error.
			if _, err := io.Copy(w, r.Body); err != nil {
				// In the case of an error in this copying process, return a server error
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
		} else {
			// In all other cases, just say hello
			w.Write([]byte("<em>Hello, world</em>"))
		}
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	})

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})

	http.ListenAndServe(":8080", nil)
}
