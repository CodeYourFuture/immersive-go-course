package main

import (
	"errors"
	"reflect"
	"testing"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func TestGrayscaleMockError(t *testing.T) {
	c := &Converter{
		cmd: func(args []string) (*imagick.ImageCommandResult, error) {
			return nil, errors.New("not implemented")
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGrayscaleMockCall(t *testing.T) {
	var args []string
	expected := []string{"convert", "input.jpg", "-set", "colorspace", "Gray", "output.jpg"}
	c := &Converter{
		cmd: func(a []string) (*imagick.ImageCommandResult, error) {
			args = a
			return &imagick.ImageCommandResult{
				Info: nil,
				Meta: "",
			}, nil
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf("incorrect arguments: expected %v, got %v", expected, args)
	}
}
