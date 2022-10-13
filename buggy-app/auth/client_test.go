package auth

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
)

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
