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
}

type authService struct {
	pb.UnimplementedAuthServer
}

func (s *authService) Verify(ctx context.Context, in *pb.Input) (*pb.Result, error) {
	return &pb.Result{
		Allow: false,
	}, nil
}

func newAuthService() *authService {
	return &authService{}
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) Run(ctx context.Context, config Config) error {
	listen := fmt.Sprintf("localhost:%d", config.Port)
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, newAuthService())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = grpcServer.Serve(lis)
	}()

	log.Printf("listening: %s", listen)

	<-ctx.Done()
	grpcServer.GracefulStop()
	wg.Wait()
	return err
}
