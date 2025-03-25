package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TestAPI struct {
   Client       *http.Client
   URL          string
   SleepTime    time.Duration
}
type RealAPI struct {
    Client      *http.Client
    URL         string
}

func (r *RealAPI) Sleep(duration time.Duration) {
    time.Sleep(duration)
}
func (r *RealAPI) Retry(w1, w2 io.Writer) {
    fmt.Fprint(w2, "\nRetrying...")
    r.DoStuff(w1, w2)
}

func (a *TestAPI) Sleep(duration time.Duration) {
   a.SleepTime = duration 
}
func (a *TestAPI) Retry(w io.Writer) {
    fmt.Fprint(w, "\nRetrying...")
}
func (r *RealAPI) DoStuff(w1, w2 io.Writer) error {
    resp, err := r.Client.Get(r.URL)
    if err != nil {
       return err 
    }
    header := resp.Header.Get("Retry-After")

    // Case1: Successfull Response
    if header == "" {
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        fmt.Fprintf(w1, "%s", body)
        return nil
    }
    // Case2: Please Wait
    var delayTime int
    calculateDelayTime(header)

    if ShouldWeWait(delayTime) {
        fmt.Fprintf(w2, "Please wait %d seconds", delayTime)
        sleepTime := time.Duration(delayTime) * time.Second

        r.Sleep(sleepTime)
        fmt.Fprint(w2, "\nRetrying...\n")
        r.DoStuff(w1, w2)
    }
    if !ShouldWeWait(delayTime) {
        fmt.Fprintf(w2, "I can't give you the weather")
        return nil
    }
    return nil


}

// Testing Purposes
// w1: writter to Stdout
// w2: writter to Stderr
func (t *TestAPI) DoStuff(w io.Writer) error {
    resp, err := t.Client.Get(t.URL)
    if err != nil {
       return err 
    }
    header := resp.Header.Get("Retry-After")

    // Case1: Successfull Response
    if header == "" {
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        fmt.Fprintf(w, "%s", body)
        return nil
    }
    // Case2: Please Wait
    delayTime, err := calculateDelayTime(header)
    if err != nil {
        return err 
    }
    if ShouldWeWait(delayTime) {
        fmt.Fprintf(w, "Please wait %d seconds", delayTime)
        sleepTime := time.Duration(delayTime) * time.Second

        t.Sleep(sleepTime)
        fmt.Fprint(w, "\nRetrying...\n")

    }else if !ShouldWeWait(delayTime) {
        fmt.Fprintf(w, "I can't give you the weather")
        return nil
    }
    return nil

    // Case3
}

// Helper Functions
func isValidHTTPFormat(str string) (t time.Time, err error) {
    parsedTime, err := time.Parse(http.TimeFormat, str)
    if err != nil {
        return time.Time{}, err
    }
    return parsedTime, nil
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
func ShouldWeWait(header int) bool {
    if header <= 5 {
        return true 
    }
    return false
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
    api := RealAPI{Client: client, URL: "http://127.0.0.1:8080"}
    api.DoStuff(os.Stdout, os.Stderr)
}
