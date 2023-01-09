package authuserctx

import (
	"context"
)

// This package has methods for adding authenticated user's ID to a context.
// For more on this idea, see https://go.dev/blog/context

type key int

// `authenticatedIdKey“ is the context key for the user identifier.
// The 0 is arbitrary -- but if another key were added to this package, it would need
// another value.
const authenticatedIdKey key = 0

func NewAuthenticatedContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, authenticatedIdKey, id)
}

func FromAuthenticatedContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(authenticatedIdKey).(string)
	return id, ok
}
