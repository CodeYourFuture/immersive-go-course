package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}

func main() {
	// Accept --input and --output arguments for the images
	inputFilepath := flag.String("input", "", "A path to an image to be processed")
	outputFilepath := flag.String("output", "", "A path to where the processed image should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputFilepath == "" || *outputFilepath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion!
	err := c.Grayscale(*inputFilepath, *outputFilepath)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	// Log what we did
	log.Printf("processed: %q to %q\n", *inputFilepath, *outputFilepath)
}
