package client

import (
    "testing"
    "bytes"
	"net/http"
    "time"
    "net/http/httptest"	
)


func TestDelayedResponse(t *testing.T) {
    t.Run("evaluate retrying given seconds", func(t *testing.T) {
        cases := []struct {
            Name                string
            Server              *httptest.Server
            SleepTime           time.Duration
            ExpectedResponse    string
        }{
            {   
                "3 seconds delayed server",
                makeDelayedServer("3"), 
                3 * time.Second,
               `Please wait 3 seconds
Retrying...
`,
            },
            {
                "5 seconds delayed server",
                makeDelayedServer("5"),
                5 * time.Second,
                `Please wait 5 seconds
Retrying...
`,
            },
            {
                "8 seconds delayed server",
                makeDelayedServer("8"),
                0 * time.Second,
                `I can't give you the weather`,
            },
        }
        for _, test := range cases {
            t.Run(test.Name, func(t *testing.T){

                api := BaseAPI{
                    Client: test.Server.Client(), 
                    URL: test.Server.URL,
                    Testing: true,
                }

                buf := bytes.Buffer{}

                if err := api.DoStuff(&buf, &buf); err != nil {
                    t.Error("something wrong happened")
                }

                // Evaluate sleep time
                if api.TestAPI.SleepTime != test.SleepTime {
                    t.Errorf("expected sleep time to be %v got %v", test.SleepTime, api.TestAPI.SleepTime)
                }
                // Evaluate buffer
                got := buf.String()
                if got != test.ExpectedResponse {
                    t.Errorf("expected response %s got %s", test.ExpectedResponse, got)
                }
            })
        }
    })
    t.Run("evaluate retrying given timestamp", func(t *testing.T) {
        cases := []struct {
            Name                string 
            Server              *httptest.Server
            SleepTime           time.Duration
            ExpectedResponse    string
        }{
            {
                "3 seconds delayed server",
                makeDelayedTimestampServer(3),
                3 * time.Second,
                `Please wait 3 seconds
Retrying...
`,
            },
            {
                "5 seconds delayed server",
                makeDelayedTimestampServer(5),
                5 * time.Second,
                `Please wait 5 seconds
Retrying...
`,
            },
            {
                "8 seconds delayed server",
                makeDelayedTimestampServer(8),
                0 * time.Second,
                `I can't give you the weather`,
            },
        }

        for _, test := range cases {
            t.Run(test.Name, func(t *testing.T) {
                api := BaseAPI{
                    Client: test.Server.Client(), 
                    URL: test.Server.URL,
                    Testing: true,
                }

                buf := bytes.Buffer{}

                if err := api.DoStuff(&buf, &buf); err != nil {
                    t.Error("something wrong happened")
                }

                // Evaluate sleep time
                if api.TestAPI.SleepTime != test.SleepTime {
                    t.Errorf("expected sleep time to be %v got %v", test.SleepTime, api.TestAPI.SleepTime)
                }
                // Evaluate buffer
                got := buf.String()
                if got != test.ExpectedResponse {
                    t.Errorf("expected response %s got %s", test.ExpectedResponse, got) 
                }
            })
        }
    })
    t.Run("evaluate retrying given undetermined server", func(t *testing.T) {
        server := makeUndeterminedServer()
        defer server.Close()

        api := BaseAPI{
            Client: server.Client(),
            URL: server.URL,
            Testing: true,
        }
        
        buf := bytes.Buffer{}
        if err := api.DoStuff(&buf, &buf); err != nil {
            t.Error("something wrong happened")
        }

        got := buf.String()
        want := `Max Time Waiting: 10s. Please do not leave
Retrying...
`

        if got != want {
            t.Errorf("got %s wanted %s", got, want)
        }
    })
}

func makeDelayedServer(seconds string) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Retry-After", seconds)
        w.WriteHeader(429)    
    }))
}
func makeUndeterminedServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Retry-After", "a while")
        w.WriteHeader(429)    
    }))
}
func makeDelayedTimestampServer(seconds int) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        time := time.Now().UTC().Add(time.Duration(seconds) * time.Second)
        timeFormatted := time.Format(http.TimeFormat)

        w.Header().Set("Retry-After", timeFormatted)
        w.WriteHeader(429)
    }))
}
