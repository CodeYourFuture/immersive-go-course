package api

import (
	"context"
	"log"
	"net/http"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util/authuserctx"
)

func (as *Service) wrapAuth(handler http.HandlerFunc) http.HandlerFunc {
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
		result, err := as.authClient.Verify(ctx, id, passwd)
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

		ctx = authuserctx.NewAuthenticatedContext(ctx, id)
		handler(w, r.WithContext(ctx))
	}
}
