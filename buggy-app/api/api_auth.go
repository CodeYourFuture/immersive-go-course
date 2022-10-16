package api

import (
	"context"
	"log"
	"net/http"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authuserctx"
)

// wrapAuth takes a handler function (likely to be the API endpoint) and wraps it with an authentication
// check using an AuthClient.
//
// If the authentication passes, it adds the authenticated user ID to the context using the authuserctx
// package, and then calls the inner handler. The ID can be retrieved later using the
// `FromAuthenticatedContext` function.
func (as *Service) wrapAuth(client auth.Client, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		id, passwd, ok := r.BasicAuth()
		// Malformed basic auth is not OK
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Use the auth client to check if this id/password combo is approved
		result, err := client.Verify(ctx, id, passwd)
		if err != nil {
			log.Printf("api: verify error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Unless we get an Allow, say no
		if result.State != auth.StateAllow {
			log.Printf("api: verify denied: id %v\n", id)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Add the ID to the context and call the inner handler
		ctx = authuserctx.NewAuthenticatedContext(ctx, id)
		handler(w, r.WithContext(ctx))
	}
}
