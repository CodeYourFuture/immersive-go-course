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
