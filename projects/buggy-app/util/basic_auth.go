package util

import (
	"encoding/base64"
	"fmt"
)

// Generate a value for a Basic Auth header:
//
//	Authorization: Basic <...>
func BasicAuthValue(id, password string) string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", id, password)))
}

// Generate a Basic auth Header value for an Authorization header
//
//	Authorization: <...>
func BasicAuthHeaderValue(id, password string) string {
	return fmt.Sprintf("Basic %s", BasicAuthValue(id, password))
}
