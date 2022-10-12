package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	passwordFile := os.Getenv("POSTGRES_PASSWORD_FILE")
	if passwordFile == "" {
		log.Fatalln("please set POSTGRES_PASSWORD_FILE environment variable")
	}

	pwdFile, err := os.ReadFile(passwordFile)
	if err != nil {
		log.Fatal(err)
	}
	passwd := string(pwdFile)

	contents, err := readDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	if len(contents) == 0 {
		log.Fatalln("path is empty directory")
	}

	for _, entry := range contents {
		if !entry.IsDir() {
			continue
		}
		dir := fmt.Sprintf("file://%s/%s", *path, entry.Name())
		url := fmt.Sprintf("postgres://postgres:%s@postgres:5432/%s?sslmode=disable", passwd, entry.Name())
		log.Printf("migrating: %q into %q database", dir, entry.Name())
		m, err := migrate.New(dir, url)
		if err != nil {
			log.Fatal(err)
		}
		m.Up()
	}

}
