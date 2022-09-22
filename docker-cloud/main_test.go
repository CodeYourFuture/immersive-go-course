package main_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker daemon: %s", err)
	}

	resource, err = pool.BuildAndRun("docker-cloud-test", "./Dockerfile", nil)
	if err != nil {
		log.Fatalf("Could not build/run: %s", err)
	}

	exitCode := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not cleanup: %s", err)
	}

	os.Exit(exitCode)
}

func TestPing(t *testing.T) {
	var err error

	var resp *http.Response
	err = pool.Retry(func() error {
		var err error
		t.Log(resource.GetPort("80/tcp"))
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("80/tcp"), "/ping"))
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatalf("response error after retries: %s", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	want := "pong"
	if got := string(body); want != got {
		t.Fatalf("bad response: want %q got %q", want, got)
	}
}
