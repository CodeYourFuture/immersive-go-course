package auth

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
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

type AuthService struct {
	util.Service

	service *authService
}

func NewAuthService() *AuthService {
	return &AuthService{
		service: newAuthService(),
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
func (as *AuthService) Run(ctx context.Context, config Config) error {
	// Connect to the database via a "pool" of connections, allowing concurrency
	pool, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	defer pool.Close()
	// Add the pool to the "inner" auth service which implements the gRPC interface
	// and responds to RPCs
	as.service.pool = pool

	// Create a TCP listener for the gRPC server to use
	listen := fmt.Sprintf("localhost:%d", config.Port)
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Set up and register the server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, as.service)

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

	// Pool is a reference to the database that we can use for queries
	pool *pgxpool.Pool
}

type userRow struct {
	id       string
	password string
	status   int
}

// Verify checks a Input for authentication validity
func (as *authService) Verify(ctx context.Context, in *pb.Input) (*pb.Result, error) {
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
		// ... either way, deny!
		return &pb.Result{
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
		return &pb.Result{
			State: pb.State_DENY,
		}, nil
	}

	// No errors from the query or the password comparison
	return &pb.Result{
		State: pb.State_ALLOW,
	}, nil
}

func newAuthService() *authService {
	return &authService{}
}
