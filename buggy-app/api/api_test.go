package api

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
)

var defaultConfig Config = Config{
	Port:           8090,
	Log:            log.Default(),
	AuthServiceUrl: "auth:8080",
}

func TestRun(t *testing.T) {
	as := New(defaultConfig)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx)
	}()

	<-time.After(1000 * time.Millisecond)
	cancel()

	wg.Wait()
	if runErr != http.ErrServerClosed {
		t.Fatal(runErr)
	}
}

func TestSimpleRequest(t *testing.T) {
	as := New(defaultConfig)

	var runErr error
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		runErr = as.Run(ctx)
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

func TestMyNotesAuthFail(t *testing.T) {
	as := New(defaultConfig)
	as.authClient = auth.NewMockClient(auth.VerifyResult{
		State: auth.StateDeny,
	})

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(as.handleMyNotes)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestMyNotesAuthFailWithAuth(t *testing.T) {
	as := New(defaultConfig)
	as.authClient = auth.NewMockClient(auth.VerifyResult{
		State: auth.StateDeny,
	})

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic ZXhhbXBsZTpleGFtcGxl")
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(as.handleMyNotes)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestMyNotesAuthFailMalformedAuth(t *testing.T) {
	as := New(defaultConfig)
	as.authClient = auth.NewMockClient(auth.VerifyResult{
		State: auth.StateDeny,
	})

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic nope")
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(as.handleMyNotes)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestMyNotesAuthPass(t *testing.T) {
	as := New(defaultConfig)
	as.authClient = auth.NewMockClient(auth.VerifyResult{
		State: auth.StateAllow,
	})

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic ZXhhbXBsZTpleGFtcGxl")
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(as.handleMyNotes)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}
