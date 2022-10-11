package auth

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"google.golang.org/grpc"
)

type Config struct {
	Port int
	Log  log.Logger
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Run starts the underlying gRPC server according to the supplied Config
// It uses the supplied context cancel signal to trigger graceful shutdown:
//
//	as := auth.NewAuthService()
//	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
//	defer stop()
//	if err := as.Run(ctx, auth.Config{
//		Port: *port,
//	}); err != nil {
//		log.Fatal(err)
//	}
func (as *AuthService) Run(ctx context.Context, config Config) error {
	listen := fmt.Sprintf("localhost:%d", config.Port)

	// Create a TCP listener for the gRPC server to use
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, newAuthService())

	// Serve on the supplied listener
	// This call blocks, so we put it in a goroutine
	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = grpcServer.Serve(lis)
	}()

	config.Log.Printf("auth service: listening: %s", listen)

	// Wait for the context cancel (e.g. from interrupt signal) before
	// gracefully shutting down any ongoing RPCs
	<-ctx.Done()
	grpcServer.GracefulStop()

	// Ensure the Serve goroutine is finished
	wg.Wait()
	return runErr
}

// Internal authService struct that implements the gRPC server interface
type authService struct {
	pb.UnimplementedAuthServer
}

// Verify checks a Input for authentication validity
func (s *authService) Verify(ctx context.Context, in *pb.Input) (*pb.Result, error) {
	return &pb.Result{
		State: pb.State_DENY,
	}, nil
}

func newAuthService() *authService {
	return &authService{}
}
