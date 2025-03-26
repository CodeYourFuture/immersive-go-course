package client

import (
    "testing"
    "bytes"
    "fmt"
	"net/http"
    "net/http/httptest"	
)

func TestSuccessfullResponse(t *testing.T) {
    t.Run("evaluate successfull response.", func(t *testing.T) {
        server := makeSuccessfullServer()
        defer server.Close()

        api := BaseAPI {
            Client: server.Client(), 
            URL: server.URL, 
            Testing: true,
        }

        buf := bytes.Buffer{}
        
        if err := api.DoStuff(&buf, &buf); err != nil {
            t.Error("something went wrong")
        }

        got := buf.String()
        want := "testing weather response"

        if got != want {
            t.Errorf("got %s wanted %s", got, want)
        }
    })
}

func makeSuccessfullServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "testing weather response")
    }))
}
