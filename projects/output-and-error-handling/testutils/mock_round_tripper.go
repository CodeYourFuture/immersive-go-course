package testutils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type responseOrError struct {
	response *http.Response
	err      error
}

type MockRoundTripper struct {
	t            *testing.T
	responses    []responseOrError
	requestCount int
}

func NewMockRoundTripper(t *testing.T) *MockRoundTripper {
	return &MockRoundTripper{
		t: t,
	}
}

func (m *MockRoundTripper) StubResponse(statusCode int, body string, header http.Header) {
	// We need to stub out a fair bit of the HTTP response in for the Go HTTP client to accept our response.
	response := &http.Response{
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        header,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Status:        fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode:    statusCode,
	}
	m.responses = append(m.responses, responseOrError{response: response})
}

func (m *MockRoundTripper) StubErrorResponse(err error) {
	m.responses = append(m.responses, responseOrError{err: err})
}

func (m *MockRoundTripper) AssertGotRightCalls() {
	m.t.Helper()

	require.Equalf(m.t, len(m.responses), m.requestCount, "Expected %d requests, got %d", len(m.responses), m.requestCount)
}

func (m *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	m.t.Helper()

	if m.requestCount >= len(m.responses) {
		m.t.Fatalf("MockRoundTripper expected %d requests but got more", len(m.responses))
	}
	resp := m.responses[m.requestCount]
	m.requestCount += 1
	return resp.response, resp.err
}
