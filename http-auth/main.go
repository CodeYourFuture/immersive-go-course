package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/time/rate"
)

func authOk(user string, pass string) bool {
	return user == os.Getenv("AUTH_USERNAME") && pass == os.Getenv("AUTH_PASSWORD")
}

// Take a rate.Limiter instance and a http.HandlerFunc and return another http.HandlerFunc that
// checks the rate limiter using `Allow()` before calling the supplied handler. If the request
// is not allowed by the limiter, a `503 Service Unavailable` Error is returned.
func rateLimit(limiter *rate.Limiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	})
}

// writeStartOfHTML is a function we can call from both the POST and GET path to start off the HTML response.
func writeStartOfHTML(w http.ResponseWriter) {
	// Indicate that we are sending back HTML
	w.Header().Add("Content-Type", "text/html")
	// Write the doctype and opening tag
	w.Write([]byte("<!DOCTYPE html>\n<html>\n"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the request is POSTing data, return what they sent back
		if r.Method == "POST" {
			// The request (r) body is an io.Reader so we can copy it into the
			// string builder and handle errors
			body := new(strings.Builder)
			if _, err := io.Copy(body, r.Body); err != nil {
				// In the case of an error in this copying process, return a server error
				log.Printf("Error copying request body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}
			writeStartOfHTML(w)
			// Write the body back to the requester in a safe way
			w.Write([]byte(html.EscapeString(body.String())))
		} else {
			writeStartOfHTML(w)
			// In all other cases, just say hello
			w.Write([]byte("<em>Hello, world</em>\n"))
			w.Write([]byte("<p>Query parameters:\n<ul>\n"))
			// Query parameters are available as a Values map[string][]string
			// https://pkg.go.dev/net/url#Values
			for k, vs := range r.URL.Query() {
				// As we're sending the query parameters straight back, we need to escape them.
				// Each value is a list, supporting query params like ?color=red&color=blue
				// so we need to iterate through each query parameter value and escape the string
				escapedVs := make([]string, 0, len(vs))
				for _, v := range vs {
					escapedVs = append(escapedVs, html.EscapeString(v))
				}
				// We can now write a list item, escaping the key and printing the escaped values list
				w.Write([]byte(fmt.Sprintf("<li>%s: [%s]</li>\n", html.EscapeString(k), strings.Join(escapedVs, ", "))))
			}
			w.Write([]byte("</ul>"))

		}
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	})

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !authOk(username, password) {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"localhost\", charset=\"UTF-8\"")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("<!DOCTYPE html>\n<html>\n"))
			w.Write([]byte(fmt.Sprintf("Hello %s!", html.EscapeString(username))))
		}
	})

	limiter := rate.NewLimiter(100, 30)

	// This endpoint is rate limited by `limiter`. The handler function is wrapped by `rateLimit`,
	// which will call it if the request is allowed under the rate limit, or automatically return
	// a 503.
	http.HandleFunc("/limited", rateLimit(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<!DOCTYPE html>\n<html>\nHello world!"))
	}))

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
