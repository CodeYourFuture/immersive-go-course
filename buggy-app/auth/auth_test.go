package auth

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRun(t *testing.T) {
	config := Config{
		Port: 8010,
		Log:  *log.Default(),
	}
	as := NewAuthService()

	var err error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = as.Run(ctx, config)
	}()

	<-time.After(1000 * time.Millisecond)
	cancel()

	wg.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleVerify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := Config{
		Port: 8010,
		Log:  *log.Default(),
	}
	as := NewAuthService()

	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx, config)
	}()

	conn, err := grpc.Dial("localhost:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthClient(conn)

	result, err := client.Verify(context.Background(), &pb.Input{})
	if err != nil {
		t.Fatalf("fail to dial: %v", err)
	}
	if result.State != pb.State_DENY {
		t.Fatalf("failed to verify, expected State_DENY, got %v", result.State)
	}

	cancel()
	wg.Wait()
	if runErr != nil {
		t.Fatalf("runErr: %v", err)
	}
}
