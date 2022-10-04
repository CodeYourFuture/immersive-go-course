package main

import (
	"bytes"
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

type Config struct {
	AwsRoleArn string
	AwsRegion  string
	S3Bucket   string
}

type Row struct {
	index          int
	url            string
	inputFilepath  string
	outputFilepath string
	outputKey      string
	outputUrl      string
}

func readAndValidateCsv(in io.Reader) ([][]string, error) {
	r := csv.NewReader(in)
	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	if len(records) <= 1 {
		return [][]string{}, fmt.Errorf("empty csv")
	}

	headerRow := records[0]
	if len(headerRow) == 0 || headerRow[0] != "url" {
		return [][]string{}, fmt.Errorf("incorrect column name: expected \"url\", got %q", headerRow[0])
	}

	return records, nil
}

func (row Row) handleRow(svc *s3.S3, config *Config) error {
	i, url, inputFilepath, outputFilepath := row.index, row.url, row.inputFilepath, row.outputFilepath
	// Create a new file that we will write to
	inputFile, err := os.Create(inputFilepath)
	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}
	defer inputFile.Close()

	// Get it from the internet!
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}
	defer res.Body.Close()

	// Ensure we got success from the server
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: download failed: row %d (%q): %s", i, url, res.Status)
	}

	// Copy the body of the response to the created file
	_, err = io.Copy(inputFile, res.Body)
	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}

	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err = imagick.ConvertImageCommand([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}

	log.Printf("processed: row %d (%q) to %q\n", i, url, outputFilepath)

	outputFile, err := os.Open(outputFilepath)
	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.S3Bucket),
		Key:    aws.String(row.outputKey),
		Body:   outputFile,
	})

	if err != nil {
		return fmt.Errorf("error: row %d (%q): %v", i, url, err)
	}

	return nil
}

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

	awsRoleArn := os.Getenv("AWS_ROLE_ARN")
	if awsRoleArn == "" {
		log.Fatalln("Please set AWS_ROLE_ARN environment variable")
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatalln("Please set AWS_REGION environment variable")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatalln("Please set S3_BUCKET environment variable")
	}

	config := &Config{
		AwsRoleArn: awsRoleArn,
		AwsRegion:  awsRegion,
		S3Bucket:   s3Bucket,
	}

	// Set up S3 session
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials.
	sess := session.Must(session.NewSession())

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the ARN.
	creds := stscreds.NewCredentials(sess, config.AwsRoleArn)

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds})

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Open the file supplied
	in, err := os.Open(*inputCsv)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file using the encoding/csv package
	inputRecords, err := readAndValidateCsv(in)
	if err != nil {
		log.Fatal(err)
	}

	outputRecords := make([][]string, 0, len(inputRecords)-1)
	outputRecords = append(outputRecords, []string{"url", "input", "output", "s3url"})

	for i, row := range inputRecords[1:] {
		url := row[0]

		prefix := fmt.Sprintf("/tmp/%d-%d", time.Now().UnixMilli(), rand.Int())
		inputFilepath := fmt.Sprintf("%s.%s", prefix, "jpg")
		outputFilepath := fmt.Sprintf("%s-out.%s", prefix, "jpg")
		// Upload just using the final part of the output filepath
		outputKey := filepath.Base(outputFilepath)
		outputUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.S3Bucket, config.AwsRegion, outputKey)

		log.Printf("downloading: row %d (%q) to %q\n", i, url, inputFilepath)

		row := Row{
			index:          i,
			url:            url,
			inputFilepath:  inputFilepath,
			outputFilepath: outputFilepath,
			outputKey:      outputKey,
			outputUrl:      outputUrl,
		}

		err := row.handleRow(svc, config)
		if err != nil {
			log.Printf("error: row %d (%q): %v", i, url, err)
			continue
		}

		outputRecords = append(outputRecords, []string{row.url, row.inputFilepath, row.outputFilepath, row.outputUrl})

		log.Printf("uploaded: row %d (%q) to %s\n", i, url, outputUrl)
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
	log.Printf("summary: %d of %d uploaded", len(outputRecords), len(inputRecords))
}
