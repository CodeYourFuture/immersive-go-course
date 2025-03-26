package main

import (
	"fmt"
	"io"
	"net/http"
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

func (a *TestAPI) Sleep(duration time.Duration) {
   a.SleepTime = duration 
}
//func (a *TestAPI) Retry(w io.Writer) {
//    fmt.Fprint(w, "\nRetrying...")
//}
func (b *BaseAPI) DoStuff(w1, w2 io.Writer) error {
    resp, err := b.Client.Get(b.URL)
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
    delayTime, err := calculateDelayTime(header); if err != nil {
        return err
    }

    if ShouldWeWait(delayTime) {
        fmt.Fprintf(w2, "Please wait %d seconds", delayTime)
        sleepTime := time.Duration(delayTime) * time.Second

        b.Sleep(sleepTime)
        b.Retry(w1, w2)
    }
    if !ShouldWeWait(delayTime) {
        fmt.Fprintf(w2, "I can't give you the weather")
        return nil
    }
    return nil


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
func main() {}
