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
	expected := ``
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}

func TestRootCmdWithFile(t *testing.T) {
	cmd := NewRoodCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"../assets/dew.txt"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	expected := `“A World of Dew” by Kobayashi Issa

A world of dew,
And within every dewdrop
A world of struggle.`
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}
