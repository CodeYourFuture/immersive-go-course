package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestRootCmd(t *testing.T) {
	cmd := NewRoodCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	expected := `root.go
root_test.go
`
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}
