package model

import (
	"reflect"
	"testing"
)

func TestTags(t *testing.T) {
	text := "This is an example #tag1 #tag2"
	expected := []string{"tag1", "tag2"}

	tags := extractTags(text)

	if !reflect.DeepEqual(expected, tags) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}

func TestTagsTrim(t *testing.T) {
	text := "This is an example #tag1    #tag2    "
	expected := []string{"tag1", "tag2"}

	tags := extractTags(text)

	if !reflect.DeepEqual(expected, tags) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}
