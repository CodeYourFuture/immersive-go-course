package fetcher

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/CodeYourFuture/immersive-go-course/projects/output-and-error-handling/testutils"
	"github.com/stretchr/testify/require"
)

func TestFetch200(t *testing.T) {
	transport := testutils.NewMockRoundTripper(t)
	transport.StubResponse(http.StatusOK, "some weather", nil)
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	weather, err := f.Fetch("http://doesnotexist.com/")

	require.NoError(t, err)
	require.Equal(t, "some weather", weather)
}

func TestFetch429(t *testing.T) {
	transport := testutils.NewMockRoundTripper(t)
	headers := make(http.Header)
	headers.Set("retry-after", "1")
	transport.StubResponse(http.StatusTooManyRequests, "server overloaded", headers)
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	start := time.Now()
	_, err := f.Fetch("http://doesnotexist.com/")
	elapsed := time.Since(start)

	require.Equal(t, ErrRetry, err)
	require.GreaterOrEqual(t, elapsed, 1*time.Second)
}

func Test500(t *testing.T) {
	transport := testutils.NewMockRoundTripper(t)
	transport.StubResponse(http.StatusInternalServerError, "Something went wrong", nil)
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	_, err := f.Fetch("http://doesnotexist.com/")

	require.EqualError(t, err, "unexpected response from server: Status code: 500 Internal Server Error, Body: Something went wrong")
}

func TestDisconnect(t *testing.T) {
	transport := testutils.NewMockRoundTripper(t)
	transport.StubErrorResponse(fmt.Errorf("connection failed"))
	defer transport.AssertGotRightCalls()

	f := &WeatherFetcher{client: http.Client{
		Transport: transport,
	}}

	_, err := f.Fetch("http://doesnotexist.com/")

	require.EqualError(t, err, "couldn't make HTTP request: Get \"http://doesnotexist.com/\": connection failed")
}

func TestParseDelay(t *testing.T) {
	// Generally when testing time, we'd inject a controllable clock rather than really using time.Now().
	futureTime := time.Date(2051, time.February, 1, 14, 00, 01, 0, time.UTC)
	futureTimeString := "Wed, 01 Feb 2051 14:00:01 GMT"

	for name, tc := range map[string]struct {
		header string
		delay  time.Duration
		err    string
	}{
		"integer seconds": {
			header: "10",
			delay:  10 * time.Second,
		},
		"decimal seconds": {
			header: "10.1",
			err:    "couldn't parse retry-after header as an integer number of seconds or a date. Value was: \"10.1\"",
		},
		"far future date:": {
			header: futureTimeString,
			delay:  time.Until(futureTime),
		},
		"empty string": {
			header: "",
			err:    `couldn't parse retry-after header as an integer number of seconds or a date. Value was: ""`,
		},
		"some text": {
			header: "beep boop",
			err:    `couldn't parse retry-after header as an integer number of seconds or a date. Value was: "beep boop"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			delay, err := parseDelay(tc.header)
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.InDelta(t, tc.delay/time.Second, delay/time.Second, 1)
			}
		})
	}
}
