package interfaces

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFilteringPipe(t *testing.T) {
    cases := []struct {
        Input           []byte
        ExpectedOutput  []byte
    } {
        {[]byte("start=1, end=10"), []byte("start=, end=")},

        {[]byte("give me the 1234"), []byte("give me the")},

        {[]byte("50 reasons"), []byte(" reasons")},
    }

    for _, test := range cases {
        buf := bytes.Buffer{}
        filteringPipe := FilteringPipe{Writer: &buf}
        filteringPipe.Write(test.Input)

        got := buf.Bytes()
        want := test.ExpectedOutput

        if reflect.DeepEqual(got, want) {
            t.Errorf("got %v wanted %v", got, want)
        }
    }
}
