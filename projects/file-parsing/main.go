package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/binary"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/csv"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/json"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/repeated_json"
)

func main() {
	format := flag.String("format", "", "Format the file is serialised in. Accepted values: json,repeated-json,csv,binary")
	file := flag.String("file", "", "Path to the file to read data from")
	flag.Parse()

	var parser parsers.Parser
	switch *format {
	case "json":
		parser = &json.Parser{}
	case "repeated-json":
		parser = &repeated_json.Parser{}
	case "csv":
		parser = &csv.Parser{}
	case "bin":
		parser = &binary.Parser{}
	case "":
		log.Fatal("format is a required argument")
	default:
		log.Fatalf("Didn't know how to parse format %q", *format)
	}

	if *file == "" {
		log.Fatal("file is a required argument")
	}
	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", *file, err)
	}
	defer f.Close()

	records, err := parser.Parse(f)
	if err != nil {
		log.Fatalf("Failed to parse file %s as %s: %v", *file, *format, err)
	}

	if len(records) == 0 {
		log.Fatal("No scores were found")
	}

	lowScore := parsers.ScoreRecord{
		HighScore: math.MaxInt32,
	}
	highScore := parsers.ScoreRecord{
		HighScore: math.MinInt32,
	}

	for _, record := range records {
		if record.HighScore > highScore.HighScore {
			highScore = record
		}
		if lowScore.HighScore < lowScore.HighScore {
			lowScore = record
		}
	}
	fmt.Printf("High score: %d from %s - congratulations!\n", highScore.HighScore, highScore.Name)
	fmt.Printf("Low score: %d from %s - commiserations!\n", lowScore.HighScore, lowScore.Name)
}
