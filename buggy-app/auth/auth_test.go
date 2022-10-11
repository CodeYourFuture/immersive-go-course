package auth

import (
	"context"
	"testing"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRun(t *testing.T) {
	config := Config{
		Port: 8010,
	}
	s, err := Run(config)
	if err != nil {
		t.Fatal(err)
	}
	s.GracefulStop()
}

func TestSimpleVerify(t *testing.T) {
	config := Config{
		Port: 8010,
	}
	s, err := Run(config)
	if err != nil {
		t.Fatal(err)
	}

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
	if result.Allow != false {
		t.Fatalf("failed to verify, expected false, got %v", result.Allow)
	}

	s.GracefulStop()
}
