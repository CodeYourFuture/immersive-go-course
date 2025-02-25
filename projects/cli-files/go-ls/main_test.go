package main

import (
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {

	t.Run("when no directory is provided", func(t *testing.T) {
		got := captureOut([]string{"go-ls"})

		want := "assets cmd go.mod main.go main_test.go \n"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("works with . ", func(t *testing.T) {
		got := captureOut([]string{"go-ls", "."})
		want := "assets cmd go.mod main.go main_test.go \n"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("works with .. ", func(t *testing.T) {
		got := captureOut([]string{"go-ls", ".."})
		want := "README.md go-cat go-ls \n"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("when a file is provided returns the file ", func(t *testing.T) {
		got := captureOut([]string{"go-ls", "go.mod"})
		want := "go.mod\n"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("when a wrong directory or file path is provided ", func(t *testing.T) {
		got := captureOut([]string{"go-ls", "mark"})
		want := "stat mark: no such file or directory\n"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

func captureOut(args []string) string {

	// save copy of std out
	old := os.Stdout

	// create read and write pipe
	r, w, _ := os.Pipe()

	// set the stdout to the pipe
	os.Stdout = w

	// Set args to test listing current directory
	os.Args = args

	// execute our function
	main()

	// close the resource
	w.Close()

	//reset the stdout back too the original
	os.Stdout = old

	// read from the output we created
	out, _ := io.ReadAll(r)

	return string(out)

}
