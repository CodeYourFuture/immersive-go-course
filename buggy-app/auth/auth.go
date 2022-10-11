package auth

import (
	"context"
	"fmt"
	"net"

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

func Run(config Config) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, newAuthService())
	go grpcServer.Serve(lis)
	return grpcServer, nil
}
