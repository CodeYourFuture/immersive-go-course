package auth

import (
	"fmt"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/cache"
	pb "github.com/CodeYourFuture/immersive-go-course/buggy-app/auth/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	Close() error
	Verify(ctx context.Context, id, passwd string) (*VerifyResult, error)
}

type VerifyResult struct {
	State string
}

var (
	StateDeny  = pb.State_name[int32(pb.State_DENY)]
	StateAllow = pb.State_name[int32(pb.State_ALLOW)]
)

// GrpcClient is meant to be used by other services to talk with the Auth service.
type GrpcClient struct {
	conn   *grpc.ClientConn
	cancel context.CancelFunc
	aC     pb.AuthClient
	cache  *cache.Cache[VerifyResult]
}

// Create a new Client for the auth service.
// Call Close() to release resources associated with this Client.
func NewClient(ctx context.Context, target string) (*GrpcClient, error) {
	return newClientWithOpts(ctx, target, defaultOpts()...)
}

// Call Close() to release resources associated with this Client.
func (c *GrpcClient) Close() error {
	// We cancel the context in case the connection is still being formed...
	c.cancel()
	// ...but according to grpc.DialContext docs, we still need to call conn.Close()
	return c.conn.Close()
}

func (c *GrpcClient) Verify(ctx context.Context, id, passwd string) (*VerifyResult, error) {
	// Check the cache to see if we have this id/passwd combo already there
	// If we do, return it so we don't contact the auth service twice
	cacheKey := c.cache.Key(fmt.Sprintf("%s:%s", id, passwd))
	if v, ok := c.cache.Get(cacheKey); ok {
		return v, nil
	}

	// Call the auth service to check the id/password we've been given
	res, err := c.aC.Verify(ctx, &pb.VerifyRequest{
		Id:       id,
		Password: passwd,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify: %w", err)
	}

	// Looking good: turn this gRPC result into our output type
	vR := &VerifyResult{
		State: pb.State_name[int32(res.State)],
	}

	// Remember this verify result for next time
	c.cache.Put(cacheKey, vR)
	return vR, nil
}

func defaultOpts() []grpc.DialOption {
	return []grpc.DialOption{
		// TODO: insecure connection should move to TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// Use this function in tests to configure the underlying client with options
func newClientWithOpts(ctx context.Context, target string, opts ...grpc.DialOption) (*GrpcClient, error) {
	// Wrapping the context WithCancel allows us to cancel the connection if the caller chooses to
	// immediately Close() the Client.
	ctx, cancel := context.WithCancel(ctx)
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &GrpcClient{
		conn:   conn,
		cancel: cancel,
		aC:     pb.NewAuthClient(conn),
		cache:  cache.New[VerifyResult](),
	}, nil
}

// Use this in tests to Mock out the client
type MockClient struct {
	result *VerifyResult
}

func NewMockClient(result *VerifyResult) *MockClient {
	return &MockClient{
		result: result,
	}
}

func (ac *MockClient) Close() error { return nil }
func (ac *MockClient) Verify(ctx context.Context, id, passwd string) (*VerifyResult, error) {
	return ac.result, nil
}
