package auth

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRun(t *testing.T) {
	passwd, err := util.ReadPasswd()
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		Port:        8010,
		DatabaseUrl: fmt.Sprintf("postgres://postgres:%s@postgres:5432/app", passwd),
		Log:         log.Default(),
	}
	as := New(config)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx)
	}()

	<-time.After(1000 * time.Millisecond)
	cancel()

	wg.Wait()
	if runErr != nil {
		t.Fatal(runErr)
	}
}

func TestSimpleVerifyDeny(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	passwd, err := util.ReadPasswd()
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		Port:        8010,
		DatabaseUrl: fmt.Sprintf("postgres://postgres:%s@postgres:5432/app", passwd),
		Log:         log.Default(),
	}
	as := New(config)

	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx)
	}()

	<-time.After(100 * time.Millisecond)

	conn, err := grpc.Dial("localhost:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)

	result, err := client.Verify(ctx, &pb.VerifyRequest{})
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("fail to dial: %v", err)
	}
	if result.State != pb.State_DENY {
		t.Fatalf("failed to verify, expected DENY, got %v", result.State)
	}

	cancel()
	wg.Wait()
	if runErr != nil {
		t.Fatalf("runErr: %v", err)
	}
}

func TestSimpleVerifyAllow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	passwd, err := util.ReadPasswd()
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		Port:        8010,
		DatabaseUrl: fmt.Sprintf("postgres://postgres:%s@postgres:5432/app", passwd),
		Log:         log.Default(),
	}
	as := New(config)

	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx)
	}()

	<-time.After(100 * time.Millisecond)

	conn, err := grpc.Dial("localhost:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)

	// Connect to DB to get a test user
	dbConn, err := pgx.Connect(ctx, config.DatabaseUrl)
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("test failed to connect: %v", err)
	}
	defer dbConn.Close(ctx)

	user := userRow{
		// banana
		password: "$2y$10$O8VPlcAPa/iKHrkdyzN1cu7TvF5Goq6nRjSdaz9uXm1zPcVgRxQnK",
		status:   "active",
	}
	err = dbConn.QueryRow(
		ctx,
		"INSERT INTO public.user (password, status) VALUES ($1, $2) RETURNING id",
		user.password,
		user.status,
	).Scan(&user.id)
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("insert failed: %v", err)
	}

	log.Printf("TestSimpleVerifyAllow: got id %s\n", user.id)

	result, err := client.Verify(ctx, &pb.VerifyRequest{
		Id:       user.id,
		Password: "banana",
	})
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("fail to verify: %v", err)
	}
	if result.State != pb.State_ALLOW {
		cancel()
		wg.Wait()
		t.Fatalf("failed to verify, expected ALLOW, got %v", result.State)
	}

	// TODO: this cleanup needs to happen regardless and be linked to the context
	_, err = dbConn.Exec(
		ctx,
		"DELETE FROM public.user WHERE id = $1",
		user.id,
	)
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatalf("failed to clean up")
	}

	cancel()
	wg.Wait()
	if runErr != nil {
		t.Fatalf("runErr: %v", err)
	}
}
