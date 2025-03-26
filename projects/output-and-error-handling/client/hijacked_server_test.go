package client

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestConnectionLost(t *testing.T) {
    server := makeHijackedServer()
    api := BaseAPI{
        Client: server.Client(),
        URL: server.URL,
        Testing: true,
    }
    
    buf := bytes.Buffer{}

    if err := api.DoStuff(&buf, &buf); err != nil {
        t.Error("something bad happened")
    }

    got := buf.String()
    want := "Connection Lost. Try again later"

    if got != want {
        t.Errorf("got %s wanted %s", got, want)
    }

}
func makeHijackedServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
    }))
}

