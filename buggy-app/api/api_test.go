package api

import (
	"context"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	config := Config{
		Log: log.Default(),
	}
	as := NewApiService()

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx, config)
	}()

	<-time.After(1000 * time.Millisecond)
	cancel()

	wg.Wait()
	if runErr != http.ErrServerClosed {
		t.Fatal(runErr)
	}
}

func TestSimpleRequest(t *testing.T) {

	config := Config{
		Port: 8090,
		Log:  log.Default(),
	}
	as := NewApiService()

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx, config)
	}()

	<-time.After(1000 * time.Millisecond)

	resp, err := http.Get("http://localhost:8090/1/my/notes.json")
	if err != nil {
		cancel()
		wg.Wait()
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		cancel()
		wg.Wait()
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	cancel()
	wg.Wait()
	if runErr != http.ErrServerClosed {
		t.Fatal(runErr)
	}
}
