package auth

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"google.golang.org/grpc"
)

// Internal grpcAuthService struct that implements the gRPC server interface
// for testing
type mockGrpcAuthService struct {
	pb.UnimplementedAuthServer

	result *pb.VerifyResponse
	err    error

	Calls int
}

func newMockGrpcService(result *pb.VerifyResponse, err error) *mockGrpcAuthService {
	return &mockGrpcAuthService{
		result: result,
		err:    err,
	}
}

// Verify checks a Input for authentication validity
func (as *mockGrpcAuthService) Verify(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	as.Calls += 1
	return as.result, as.err
}

func TestClientCreate(t *testing.T) {
	config := Config{
		Port: 8010,
		Log:  log.Default(),
	}
	as := New(config)

	var err error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = as.Run(ctx)
	}()

	client, err := newClientWithOpts(ctx, "localhost:8010", defaultOpts()...)
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(500 * time.Millisecond)
	client.Close()

	<-time.After(1000 * time.Millisecond)
	cancel()

	wg.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientError(t *testing.T) {
	opts := append(defaultOpts(), grpc.WithBlock())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := newClientWithOpts(ctx, "localhost:8010", opts...)
	if err == nil {
		t.Fatal("did not error")
	}
}

func TestClientClose(t *testing.T) {
	client, err := NewClient(context.Background(), "localhost:8010")
	if err != nil {
		t.Fatal(err)
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientVerifyDeny(t *testing.T) {
	listen := "localhost:8010"
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	pbStateExpected, stateExpected := pb.State_DENY, StateDeny

	mockService := newMockGrpcService(&pb.VerifyResponse{
		State: pbStateExpected,
	}, nil)

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, mockService)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = grpcServer.Serve(lis)
	}()

	done := func() {
		cancel()
		grpcServer.GracefulStop()
		wg.Wait()
	}

	client, err := NewClient(ctx, listen)
	if err != nil {
		done()
		t.Fatal(err)
	}

	res, err := client.Verify(ctx, "example", "example")
	if err != nil {
		done()
		t.Fatal(err)
	}

	err = client.Close()
	if err != nil {
		done()
		t.Fatal(err)
	}

	if res.State != stateExpected {
		done()
		t.Fatalf("verify state: expected %s, got %s\n", stateExpected, res.State)
	}

	done()
	if runErr != nil && runErr != grpc.ErrServerStopped {
		t.Fatal(runErr)
	}
}

func TestClientVerifyAllow(t *testing.T) {
	listen := "localhost:8010"
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	pbStateExpected, stateExpected := pb.State_ALLOW, StateAllow

	mockService := newMockGrpcService(&pb.VerifyResponse{
		State: pbStateExpected,
	}, nil)

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, mockService)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = grpcServer.Serve(lis)
	}()

	done := func() {
		cancel()
		grpcServer.GracefulStop()
		wg.Wait()
	}

	client, err := NewClient(ctx, listen)
	if err != nil {
		done()
		t.Fatal(err)
	}

	res, err := client.Verify(ctx, "example", "example")
	if err != nil {
		done()
		t.Fatal(err)
	}

	err = client.Close()
	if err != nil {
		done()
		t.Fatal(err)
	}

	if res.State != stateExpected {
		done()
		t.Fatalf("verify state: expected %s, got %s\n", stateExpected, res.State)
	}

	done()
	if runErr != nil && runErr != grpc.ErrServerStopped {
		t.Fatal(runErr)
	}
}

func TestClientVerifyCacheResult(t *testing.T) {
	listen := "localhost:8010"
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	pbStateExpected, stateExpected := pb.State_ALLOW, StateAllow

	mockService := newMockGrpcService(&pb.VerifyResponse{
		State: pbStateExpected,
	}, nil)

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, mockService)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = grpcServer.Serve(lis)
	}()

	done := func() {
		cancel()
		grpcServer.GracefulStop()
		wg.Wait()
	}

	client, err := NewClient(ctx, listen)
	if err != nil {
		done()
		t.Fatal(err)
	}

	res1, err := client.Verify(ctx, "example", "example")
	if err != nil {
		done()
		t.Fatal(err)
	}

	res2, err := client.Verify(ctx, "example", "example")
	if err != nil {
		done()
		t.Fatal(err)
	}

	err = client.Close()
	if err != nil {
		done()
		t.Fatal(err)
	}

	if res1.State != stateExpected {
		done()
		t.Fatalf("verify state in first call: expected %s, got %s\n", stateExpected, res1.State)
	}

	if res2.State != stateExpected {
		done()
		t.Fatalf("verify state in second call: expected %s, got %s\n", stateExpected, res2.State)
	}

	if mockService.Calls > 1 {
		done()
		t.Fatalf("verify did not cache result: %d calls to service, expected 1", mockService.Calls)
	}

	done()
	if runErr != nil && runErr != grpc.ErrServerStopped {
		t.Fatal(runErr)
	}
}
