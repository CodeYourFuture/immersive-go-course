package authctx

import (
	"context"
)

// This package has methods for adding authenticated user's ID to a context.

type key int

// authenticatedIdKey is the context key for the user identifier.
const authenticatedIdKey key = 0

func NewAuthenticatedContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, authenticatedIdKey, id)
}

func FromAuthenticatedContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(authenticatedIdKey).(string)
	return id, ok
}
