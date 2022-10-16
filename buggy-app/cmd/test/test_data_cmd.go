package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// This package is a CLI tool for interacting with the database to create/update/delete data for testing.

type Flags struct {
	cmd string

	hostport string
	db       string

	// User flags
	passwd string

	// Note flags
	content string
	owner   string
}

func usage() {
	flag.Usage()
	os.Exit(2)
}

func main() {
	f := &Flags{}
	baseFlags(f, flag.CommandLine)
	if len(os.Args) < 2 {
		log.Println("error: not enough arguments, please supply one of: user")
		usage()
	}

	f.cmd = os.Args[1]
	userFlagSet := userFlags(f)
	noteFlagSet := noteFlags(f)

	var err error
	switch f.cmd {
	case "user":
		err = userFlagSet.Parse(os.Args[2:])
	case "note":
		err = noteFlagSet.Parse(os.Args[2:])
	default:
		log.Println("error: command not recognised")
		usage()
	}

	if err != nil {
		log.Println("error: could not parse flags")
		usage()
	}

	if os.Getenv("POSTGRES_PASSWORD_FILE") == "" {
		os.Setenv("POSTGRES_PASSWORD_FILE", "volumes/secrets/postgres-passwd")
	}
	dbPasswd, err := util.ReadPasswdFile()
	if err != nil {
		log.Fatal(err)
	}

	// The NotifyContext will signal Done when these signals are sent, allowing others to shut down safely
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", dbPasswd, f.hostport, f.db)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	switch f.cmd {
	case "user":
		err = userCmd(ctx, f, conn)
	case "note":
		err = noteCmd(ctx, f, conn)
	default:
		log.Fatal("expected user")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func baseFlags(f *Flags, fs *flag.FlagSet) {
	fs.StringVar(&f.hostport, "hostport", "localhost:5432", "host:port of Postgres")
	fs.StringVar(&f.db, "db", "app", "target database")
}

func userFlags(f *Flags) *flag.FlagSet {
	fs := flag.NewFlagSet("user", flag.ExitOnError)
	fs.StringVar(&f.passwd, "password", "password", "password of the created user")
	return fs
}

func userCmd(ctx context.Context, f *Flags, conn *pgx.Conn) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(f.passwd), 10)
	if err != nil {
		return fmt.Errorf("user: could not hash password, %w", err)
	}

	var id string
	err = conn.QueryRow(ctx, "INSERT INTO public.user (status, password) VALUES ($1, $2) RETURNING id", 1, hash).Scan(&id)
	if err != nil {
		return fmt.Errorf("user: could not insert user, %w", err)
	}
	log.Printf("new user created\n")
	log.Printf("\tid: %s\n", id)
	log.Printf("\tpassword: %s\n", f.passwd)
	log.Printf("base64 for auth: %s\n", util.BasicAuthValue(id, f.passwd))
	return nil
}

func noteFlags(f *Flags) *flag.FlagSet {
	fs := flag.NewFlagSet("note", flag.ExitOnError)
	fs.StringVar(&f.content, "content", "Example note content", "content of the created note")
	fs.StringVar(&f.owner, "owner", "", "owner of the created note")
	return fs
}

func noteCmd(ctx context.Context, f *Flags, conn *pgx.Conn) error {
	if f.owner == "" {
		return errors.New("note: please supply an owner with -owner")
	}
	var owner string
	err := conn.QueryRow(ctx, "SELECT id FROM public.user WHERE id = $1", f.owner).Scan(&owner)
	if err != nil {
		return fmt.Errorf("note: could not find owner, %w", err)
	}

	var id string
	err = conn.QueryRow(ctx, "INSERT INTO public.note (owner, content) VALUES ($1, $2) RETURNING id", f.owner, f.content).Scan(&id)
	if err != nil {
		return fmt.Errorf("note: could not insert note, %w", err)
	}
	log.Printf("new note created\n")
	log.Printf("\tid: %s\n", id)
	log.Printf("\towner: %s\n", f.owner)
	log.Printf("\tcontent: %q\n", f.content)
	return nil
}
