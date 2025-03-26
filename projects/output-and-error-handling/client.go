package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)
type BaseAPI struct {
    Client       *http.Client
    URL         string 
    Testing     bool
    TestAPI     TestAPI
}

type TestAPI struct {
   SleepTime    time.Duration
}

func (a *TestAPI) Sleep(duration time.Duration) {
   a.SleepTime = duration 
}

func (b *BaseAPI) Sleep(duration time.Duration) {
    if b.Testing{
        b.TestAPI.Sleep(duration)
    }else {
        time.Sleep(duration)
    }
}
func (b *BaseAPI) Retry(w1, w2 io.Writer) {
    fmt.Fprint(w2, "\nRetrying...\n")
    if !b.Testing { 
        b.DoStuff(w1, w2)
    }
}

func (b *BaseAPI) DoStuff(w1, w2 io.Writer) error {
    resp, err := b.Client.Get(b.URL)

    // 1. Connection Lost
    if err != nil {
        fmt.Fprint(w2, "Connection Lost. Try again later")
        return nil
    } 

    // 2. Successfull Response
    if resp.StatusCode == 200 {
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        fmt.Fprintf(w1, "%s", body)
        return nil
    }

    header := resp.Header.Get("Retry-After")

    // 3. Delayed Response
    if resp.StatusCode == 429 {
        format := evaluateHeaderFormat(header)
        
        // 3.1 Valid Delayed Response
        if format == "seconds" || format == "timestamp" {
            delayTime, err := calculateDelayTime(header); if err != nil {
            return err
            }
            if delayTime <= 5 {
                fmt.Fprintf(w2, "Please wait %d seconds", delayTime)
                sleepTime := time.Duration(delayTime) * time.Second

                b.Sleep(sleepTime)
                b.Retry(w1, w2)
            }else {
                fmt.Fprintf(w2, "I can't give you the weather")
                return nil
            }
        } 
        // 3.2 Invalid Delayed Response
        if format == "invalid" {
            fmt.Fprint(w2, "Max Time Waiting: 10s. Please do not leave")
            sleeTime := 10 * time.Second
            b.Sleep(sleeTime)
            b.Retry(w1, w2)
        }
    }
    return nil
}

func evaluateHeaderFormat(s string) string {
   if len(s) == 1 { return "seconds" }
   if _, err := time.Parse(http.TimeFormat, s); err == nil {
        return "timestamp"
   }
   return "invalid"
}
// Helper Functions
func isValidHTTPFormat(str string) (time.Time, error) {
    timeParsed, err := time.Parse(http.TimeFormat, str)
    if err != nil {
        return time.Time{}, err
    }
    return timeParsed, nil
}

func calculateDelayTime(header string) (int, error) {
    var delayTime int
    if parsedTime, err := isValidHTTPFormat(header); err == nil {
        delayTime = int(parsedTime.Sub(time.Now()).Seconds()) + 1
        return delayTime, nil
    }

    delayTime, err := headerToInt(header)
    if err != nil {
        return 0, err 
    }
    return delayTime, nil
}

func headerToInt(header string) (int, error) {
    num, err := strconv.Atoi(header)
    if err !=  nil {
        return 0, err
    }
    return num, nil
}
func main() {
    client := &http.Client{}

    api := BaseAPI{
        Client: client,
        URL: "http://127.0.0.1:8080",
        Testing: false,
    }

    api.DoStuff(os.Stdout, os.Stderr)
}
