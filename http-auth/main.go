package main

import (
	"fmt"
	"html"
	"io"
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
func rateLimit(lim *rate.Limiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if lim.Allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Indicate that we are sending back HTML
		w.Header().Add("Content-Type", "text/html")
		// Write the doctype and opening tag regardless of method
		w.Write([]byte("<!DOCTYPE html>\n<html>\n"))
		// If the request is POSTing data, return what they sent back
		if r.Method == "POST" {
			// The request (r) body is an io.Reader so we can copy it into the
			// string builder and handle errors
			body := new(strings.Builder)
			if _, err := io.Copy(body, r.Body); err != nil {
				// In the case of an error in this copying process, return a server error
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
			}
			// Write the body back to the requester in a safe way
			w.Write([]byte(html.EscapeString(body.String())))
		} else {
			// In all other cases, just say hello
			w.Write([]byte("<em>Hello, world</em>\n"))
			w.Write([]byte("<p>Query parameters:\n<ul>\n"))
			// Query parameters are available as a Values map[string][]string
			// https://pkg.go.dev/net/url#Values
			for k, vs := range r.URL.Query() {
				// As we're sending the query parameters straight back, we need to escape them.
				// Each value is a list, supporting query params like ?color=red&color=blue
				// so we need to iterate through each query parameter value and escape the string
				escaped_vs := make([]string, len(vs))
				for i, v := range vs {
					escaped_vs[i] = html.EscapeString(v)
				}
				// We can now write a list item, escaping the key and printing the escaped values list
				// TODO: is the use of %s here unsafe? https://pkg.go.dev/fmt
				w.Write([]byte(fmt.Sprintf("<li>%s: %s</li>\n", html.EscapeString(k), escaped_vs)))
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

	lim := rate.NewLimiter(100, 30)

	// This endpoint is rate limited by `lim`. The handler function is wrapped by `rateLimit`, which
	// will call it if the request is allowed under the rate limit, or automatically return a 503
	http.HandleFunc("/limited", rateLimit(lim, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<!DOCTYPE html>\n<html>\nHello world!"))
	}))

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})

	http.ListenAndServe(":8080", nil)
}
