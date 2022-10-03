package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	// We need a file to read from...
	inputCsv := flag.String("input", "", "A path to a CSV with a `url` column, containing URLs for images to be processed")
	// ... and a file to write to
	outputCsv := flag.String("output", "", "Location that the output of this command should be written")

	flag.Parse()
	if *inputCsv == "" || *outputCsv == "" {
		flag.Usage()
		os.Exit(1)
	}

	awsRoleUrn := os.Getenv("AWS_ROLE_URN")
	if awsRoleUrn == "" {
		log.Fatalln("Please set AWS_ROLE_URN environment variable")
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRoleUrn == "" {
		log.Fatalln("Please set AWS_REGION environment variable")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if awsRoleUrn == "" {
		log.Fatalln("Please set S3_BUCKET environment variable")
	}

	// Set up S3 session
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials.
	sess := session.Must(session.NewSession())

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the ARN.
	creds := stscreds.NewCredentials(sess, awsRoleUrn)

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds})

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Open the file supplied
	in, err := os.Open(*inputCsv)
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

	outputRecords := make([][]string, 0, len(records)-1)
	outputRecords = append(outputRecords, []string{"url", "input", "output", "s3url"})

	for i, row := range records[1:] {
		url := row[0]

		inputFilepath := fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), "jpg")
		outputFilepath := fmt.Sprintf("/tmp/%d-%d.%s", time.Now().UnixMilli(), rand.Int(), "jpg")

		log.Printf("downloading: row %d (%q) to %q\n", i, url, inputFilepath)

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

		// Ensure we got success from the server
		if res.StatusCode != http.StatusOK {
			log.Fatalf("download failed: row %d (%q): %s\n", i, url, res.Status)
		}

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

		log.Printf("processed: row %d (%q) to %q\n", i, url, outputFilepath)

		outputFile, err := os.Open(outputFilepath)
		if err != nil {
			log.Fatal(err)
		}

		// Upload just using the final part of the output filepath
		s3Key := filepath.Base(outputFilepath)

		// Uploads the object to S3. The Context will interrupt the request if the
		// timeout expires.
		_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(s3Key),
			Body:   outputFile,
		})
		if err != nil {
			log.Fatalf("failed to upload object: %v\n", err)
		}

		log.Printf("uploaded: row %d (%q) to %s/%s\n", i, url, s3Bucket, s3Key)

		outputUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Bucket, awsRegion, s3Key)

		outputRecords = append(outputRecords, []string{url, inputFilepath, outputFilepath, outputUrl})
	}

	// Turn the output records into a CSV
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	err = w.WriteAll(outputRecords)
	if err != nil {
		log.Fatalf("failed to create CSV from output records: %v\n", err)
	}
	err = os.WriteFile(*outputCsv, buf.Bytes(), os.FileMode(0644))
	if err != nil {
		log.Fatalf("failed to write output records to file: %v\n", err)
	}

	log.Printf("output: %q", *outputCsv)
}
