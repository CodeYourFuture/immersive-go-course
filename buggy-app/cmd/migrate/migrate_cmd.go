package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
)

func readDir(path string) ([]os.DirEntry, error) {
	// Stat the file so we can check if it's a directory or not before
	// we try to read it as a directory using ReadDir. os.ReadDir will
	// generates an error if the thing you pass to it is not a directory.
	// https://pkg.go.dev/os#Stat
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// We can only list the contents of a directory.
	// https://pkg.go.dev/io/fs#FileInfo
	if !fileInfo.IsDir() {
		return nil, errors.New("path must be a directory")
	}

	// Read this directory to get a list of files
	// https://pkg.go.dev/os#ReadDir
	contents, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func main() {
	path := flag.String("path", "", "Path to migrations directory")
	hostport := flag.String("hostport", "postgres:5432", "Host:port of Postgres")

	flag.Parse()
	if len(flag.Args()) == 0 || *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	passwd, err := util.ReadPasswdFile()
	if err != nil {
		log.Fatal(err)
	}

	contents, err := readDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	if len(contents) == 0 {
		log.Fatalln("path is empty directory")
	}

	// TODO: this is hack to make the migration wait for Postgres to be ready. The correct
	// thing is to detect the retryable error ("pq: the database system is starting up") and
	// retry, with exponential backoff.
	<-time.After(500 * time.Millisecond)

	for _, entry := range contents {
		// We only want to migrate directories
		if !entry.IsDir() {
			continue
		}
		// Build a file:// and a postres:// URL to migrate into
		dir := fmt.Sprintf("file://%s/%s", *path, entry.Name())
		url := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", passwd, *hostport, entry.Name())
		log.Printf("migrate: %q into %q database", dir, entry.Name())
		// Prepare the migration
		m, err := migrate.New(dir, url)
		if err != nil {
			log.Fatal(err)
		}
		// Do it, according to the argument
		switch flag.Arg(0) {
		case "up":
			err = m.Up()
		case "down":
			err = m.Down()
		default:
			log.Fatal("expected one of up or down")
		}
		if err != nil {
			// The NoChange error is not a problem
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("migrate: %s: no change", dir)
			} else {
				log.Fatal(err)
			}
		}
	}

	log.Println("migrate: complete")
}
