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
)

func main() {
	// We need a file to read from
	input := flag.String("input", "", "A path to a CSV with a `url` column, containing URLs for images to be processed")
	flag.Parse()
	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}

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

		filepath := fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), "jpg")

		log.Printf("downloading: row %d (%q) to %q", i, url, filepath)

		// Create a new file that we will write to
		out, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// Get it from the internet!
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		// TODO: check status code

		// Copy the body of the response to the created file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}
