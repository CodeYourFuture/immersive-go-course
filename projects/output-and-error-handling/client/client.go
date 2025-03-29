package client

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)
// My design includes a boolean to determine testing purposes at the moment of creation,
// and a spy struct (TestAPI) to implement basic mocking.
type BaseAPI struct {
    Client       *http.Client
    URL         string 
    Testing     bool
    TestAPI     TestAPI
}

// Spy API.
type TestAPI struct {
   SleepTime    time.Duration
}
// Spy method for testing purposes.
// In this way, the tests does not really sleep.
func (a *TestAPI) Sleep(duration time.Duration) {
   a.SleepTime = duration 
}
// Checks testing purposes and decide to call the real or spy method.
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
// For real implementation -> w1: std output, w2: std error.
// For testing just pass a buffer for both.
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
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return err
        }

        fmt.Fprintf(w1, "%s", body)
        return nil
    }

    header := resp.Header.Get("Retry-After")

    // 3. Delayed Response
    if resp.StatusCode == 429 {
        format, delayTime, err := evaluateHeaderFormat(header)

        if err != nil {
            return err
        }

        // 3.1 Valid Delayed Response
        if format == "seconds" || format == "timestamp" {
            // If time is less than 5 sec retry.
            if delayTime <= 5 {
                fmt.Fprintf(w2, "Please wait %d seconds", delayTime)
                sleepTime := time.Duration(delayTime) * time.Second

                b.Sleep(sleepTime)
                b.Retry(w1, w2)
            }else {
                // Abort otherwise.
                fmt.Fprintf(w2, "I can't give you the weather")
            }
        } 
        // 3.2 Invalid Delayed Response
        if format == "invalid" {
            // If not time specified, wait 10 sec and retry.
            fmt.Fprint(w2, "Max Time Waiting: 10s. Please do not leave")
            sleeTime := 10 * time.Second

            b.Sleep(sleeTime)
            b.Retry(w1, w2)
        }
    }
    return nil
}

// Helper Functions

// Evaluates the Retry-After header.
func evaluateHeaderFormat(s string) (format string, delayTime int, err error) {
    // 1 char means header in seconds.
    if len(s) == 1 {
        if delayTime, err = calculateDelayInSeconds(s); err != nil {
            return "", 0, err
        }
        return "seconds", delayTime, nil
    }
    // HTTP Date format means timestamp.
    if delayTime, err = isValidHTTPFormat(s); err == nil {
         return "timestamp", delayTime, nil
    }
    // If not in seconds or timestamp it must be invalid.
    return "invalid", 0, nil
}

// Evaluates the delay time on an http-date Retry-After response header.
func isValidHTTPFormat(str string) (delayTime int, err error) {
    // time.Parse: parses a formatted string and returns the time value it represents.
    timeParsed, err := time.Parse(http.TimeFormat, str)
    if err != nil {
        return 0, err
    }
    // I added 1 to the result to round the substraction of the seconds between
    // the future date and the actual date to be accurate.
    delayTime = int(timeParsed.Sub(time.Now()).Seconds()) + 1 
    return delayTime, nil
}
// Evaluates the delay time of a delay-seconds Retry-After response header.
func calculateDelayInSeconds(header string) (delayTime int, err error) {
    if delayTime, err = strconv.Atoi(header); err != nil {
        return 0, err
    }
    return delayTime, nil
}
