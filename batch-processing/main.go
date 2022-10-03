package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	// We need a file to read from
	input := flag.String("input", "", "A path to a CSV with a `url` column, containing URLs for images to be processed")
	flag.Parse()
	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Open the file supplied
	in, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file using the encoding/csv package
	r := csv.NewReader(in)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: validate `records`

	for i, row := range records[1:] {
		url := row[0]

		inputFilepath := fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), "jpg")
		outputFilepath := fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), "jpg")

		log.Printf("downloading: row %d (%q) to %q", i, url, inputFilepath)

		// Create a new file that we will write to
		inputFile, err := os.Create(inputFilepath)
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()

		// Get it from the internet!
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		// TODO: check status code

		// Copy the body of the response to the created file
		_, err = io.Copy(inputFile, res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the image to grayscale using imagemagick
		// We are directly calling the convert command
		imagick.ConvertImageCommand([]string{
			"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("processed: row %d (%q) to %q", i, url, outputFilepath)
	}

}
