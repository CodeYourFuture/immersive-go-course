package main

import (
	"reflect"
	"testing"
)

func TestPipeline(t *testing.T) {
	input := "test"
	expected := UploadResult{
		ProcessResult: ProcessResult{
			DownloadResult: DownloadResult{
				Url:      "test",
				Filepath: "test",
				Error:    nil,
			},
			Filepath: "test",
			Error:    nil,
		},
		Url:   "test",
		Error: nil,
	}
	result := pipeline(input)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
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
