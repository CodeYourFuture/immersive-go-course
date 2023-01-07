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

// This package is a CLI tool for interacting with the database to create/update/delete data for testing. It
// can create users and notes.
//
// Use it like this:
//
//	> go run ./cmd/test user -password banana
//		2022/10/16 16:40:56 new user created
//		2022/10/16 16:40:56 	id: FxoAB2gl
//		2022/10/16 16:40:56 	password: banana
//		2022/10/16 16:40:56 base64 for auth: RnhvQUIyZ2w6YmFuYW5h
//
// > go run ./cmd/test note -owner FxoAB2gl
// 		2022/10/16 16:41:42 new note created
// 		2022/10/16 16:41:42 	id: 0fxh25SJ
// 		2022/10/16 16:41:42 	owner: FxoAB2gl
// 		2022/10/16 16:41:42 	content: "Example note content"
//

type Flags struct {
	cmd string

	hostport string
	db       string

	// Number of entities to generate
	n int

	// User flags
	passwd string
	status string

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
	if len(os.Args) < 2 {
		log.Println("error: not enough arguments, expected one of: user, note")
		usage()
	}

	f.cmd = os.Args[1]
	userFlagSet := userFlags(f)
	noteFlagSet := noteFlags(f)

	var err error
	var fs *flag.FlagSet
	switch f.cmd {
	case "user":
		fs = userFlagSet
	case "note":
		fs = noteFlagSet
	default:
		log.Println("error: command not recognised")
		usage()
	}

	// Attach base flags to this flag and parse
	baseFlags(f, fs)
	err = fs.Parse(os.Args[2:])
	if err != nil {
		log.Println("error: could not parse flags")
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", f.cmd)
		usage()
	}

	// Set up a default POSTGRES_PASSWORD_FILE because we know where it's likely to be...
	if os.Getenv("POSTGRES_PASSWORD_FILE") == "" {
		os.Setenv("POSTGRES_PASSWORD_FILE", "volumes/secrets/postgres-passwd")
	}
	// ... and the read it. $POSTGRES_USER will still take precedence.
	dbPasswd, err := util.ReadPasswd()
	if err != nil {
		log.Fatal(err)
	}

	// The NotifyContext will signal Done when these signals are sent, allowing others to shut down safely
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// Connect to the database
	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", dbPasswd, f.hostport, f.db)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	i := 0
	for i < f.n {
		i += 1
		// Run the right command
		switch f.cmd {
		case "user":
			err = userCmd(ctx, f, conn)
		case "note":
			err = noteCmd(ctx, f, conn)
		default:
			log.Fatalf("unrecognised command: %s", f.cmd)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Base flags apply to all commands
func baseFlags(f *Flags, fs *flag.FlagSet) {
	fs.StringVar(&f.hostport, "hostport", "localhost:5432", "host:port of Postgres")
	fs.StringVar(&f.db, "db", "app", "target database")
	fs.IntVar(&f.n, "n", 1, "number of entities to generate")
}

// Flags associated with the user command
func userFlags(f *Flags) *flag.FlagSet {
	fs := flag.NewFlagSet("user", flag.ExitOnError)
	fs.StringVar(&f.passwd, "password", "password", "password of the created user")
	fs.StringVar(&f.status, "status", "active", "status of the created user")
	return fs
}

// Create a user from command-line configuration
func userCmd(ctx context.Context, f *Flags, conn *pgx.Conn) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(f.passwd), 10)
	if err != nil {
		return fmt.Errorf("user: could not hash password, %w", err)
	}

	if f.status != "active" && f.status != "inactive" {
		return fmt.Errorf("user: invalid status, %s", f.status)
	}

	var id string
	err = conn.QueryRow(ctx, "INSERT INTO public.user (status, password) VALUES ($1, $2) RETURNING id", f.status, hash).Scan(&id)
	if err != nil {
		return fmt.Errorf("user: could not insert user, %w", err)
	}
	log.Printf("new user created\n")
	log.Printf("\tid: %s\n", id)
	log.Printf("\tstatus: %s\n", f.status)
	log.Printf("\tpassword: %s\n", f.passwd)
	log.Printf("base64 for auth: %s\n", util.BasicAuthValue(id, f.passwd))
	return nil
}

// Flags associated with the note command
func noteFlags(f *Flags) *flag.FlagSet {
	fs := flag.NewFlagSet("note", flag.ExitOnError)
	fs.StringVar(&f.content, "content", "Example note content", "content of the created note")
	fs.StringVar(&f.owner, "owner", "", "owner of the created note")
	return fs
}

// Create a note from command-line configuration
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
