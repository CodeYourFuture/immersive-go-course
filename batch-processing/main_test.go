package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPipeline(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("assets/dog.jpg")
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(w, f)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	input := ts.URL
	expected := UploadResult{
		ProcessResult: ProcessResult{
			DownloadResult: DownloadResult{
				Url:      ts.URL,
				Filepath: "test",
				Error:    nil,
			},
			Filepath: "test",
			Error:    nil,
		},
		Url:   ts.URL,
		Error: nil,
	}
	result := pipeline(input)
	if !strings.HasSuffix(result.ProcessResult.Filepath, ".jpg") {
		t.Fatalf("ProcessResult.Filepath: is not .jpg, got %s", result.ProcessResult.Filepath)
	}
	if result.ProcessResult.DownloadResult.Url != expected.ProcessResult.DownloadResult.Url {
		t.Fatalf("ProcessResult.DownloadResult.Url: expected %+v, got %+v", expected, result)
	}
}

func TestConsumeHeader(t *testing.T) {
	input := []string{"url"}
	c := make(chan string)
	err := consumeHeader(input, Context{
		Input: c,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestConsumeHeaderErrorColumnName(t *testing.T) {
	input := []string{"error"}
	expected := "CSV header is incorrect, expected \"url\", got \"error\""
	c := make(chan string)
	err := consumeHeader(input, Context{
		Input: c,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != expected {
		t.Fatalf("incorrect error: expected \"%s\", got \"%s\"", expected, err.Error())
	}
}

func TestConsumeHeaderErrorLength(t *testing.T) {
	input := []string{"url", "error"}
	expected := "CSV header is incorrect length, expected 1, got 2"
	c := make(chan string)
	err := consumeHeader(input, Context{
		Input: c,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != expected {
		t.Fatalf("incorrect error: expected \"%s\", got \"%s\"", expected, err.Error())
	}
}

func TestConsumeRow(t *testing.T) {
	input := []string{"test"}
	expected := "test"
	c := make(chan string, 1)
	err := consumeRow(input, Context{
		Input: c,
	})
	if err != nil {
		t.Fatal(err)
	}
	result := <-c
	if expected != result {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}
