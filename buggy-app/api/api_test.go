package api

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	config := Config{}
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
	if runErr != nil {
		t.Fatal(runErr)
	}
}
