package auth

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type Config struct {
	Port        int
	DatabaseUrl string
	Log         *log.Logger
}

type Service struct {
	config      Config
	grpcService *grpcAuthService
}

func New(config Config) *Service {
	return &Service{
		config:      config,
		grpcService: newGrpcService(),
	}
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
func (as *Service) Run(ctx context.Context) error {
	// Connect to the database via a "pool" of connections, allowing concurrency
	pool, err := pgxpool.New(ctx, as.config.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	defer pool.Close()
	// Add the pool to the "inner" auth service which implements the gRPC interface
	// and responds to RPCs
	as.grpcService.pool = pool

	// Create a TCP listener for the gRPC server to use
	listen := fmt.Sprintf(":%d", as.config.Port)
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, as.grpcService)

	// Serve on the supplied listener
	// This call blocks, so we put it in a goroutine
	var runErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = grpcServer.Serve(lis)
	}()

	as.config.Log.Printf("auth service: listening: %s", listen)

	// Wait for the context cancel (e.g. from interrupt signal) before
	// gracefully shutting down any ongoing RPCs
	<-ctx.Done()
	grpcServer.GracefulStop()

	// Ensure the Serve goroutine is finished
	wg.Wait()
	return runErr
}

// Internal grpcAuthService struct that implements the gRPC server interface
type grpcAuthService struct {
	pb.UnimplementedAuthServer

	// Pool is a reference to the database that we can use for queries
	pool *pgxpool.Pool
}

func newGrpcService() *grpcAuthService {
	return &grpcAuthService{}
}

type userRow struct {
	id       string
	password string
	status   string
}

// Verify checks a Input for authentication validity
func (as *grpcAuthService) Verify(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	log.Printf("verify: id %v, start\n", in.Id)

	// Look for this user in the database
	var row userRow
	err := as.pool.QueryRow(ctx,
		"SELECT id, password, status FROM public.user WHERE id = $1",
		in.Id,
	).Scan(&row.id, &row.password, &row.status)
	// Error can be no rows or a real error...
	if err != nil {
		// No rows is not an error that needs logging
		if err != pgx.ErrNoRows {
			log.Printf("verify: query error: %v\n", err)
		}
		log.Printf("verify: id %v, deny (query)\n", in.Id)
		// ... either way, deny!
		return &pb.VerifyResponse{
			State: pb.State_DENY,
		}, nil
	}

	// bcrypt require us to compare the input to the hash directly
	// https://auth0.com/blog/hashing-in-action-understanding-bcrypt/
	err = bcrypt.CompareHashAndPassword([]byte(row.password), []byte(in.Password))
	if err != nil {
		// Mismatched hash and password is OK, but other errors need logging
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Printf("verify: compare error: %v\n", err)
		}
		log.Printf("verify: id %v, deny (password)\n", in.Id)
		return &pb.VerifyResponse{
			State: pb.State_DENY,
		}, nil
	}

	log.Printf("verify: id %v, allow\n", in.Id)
	// No errors from the query or the password comparison
	return &pb.VerifyResponse{
		State: pb.State_ALLOW,
	}, nil
}
