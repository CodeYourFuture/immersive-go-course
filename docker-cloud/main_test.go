package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

func TestEndpoints(t *testing.T) {
	t.Run("testing index handler", func(t *testing.T){
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		indexHandler(response, request)
		assertStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), "Hello, world.")
	})
	t.Run("testing ping handler", func(t *testing.T){
		request, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		response := httptest.NewRecorder()
		pingHandler(response, request)
		assertStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), "Hello!")
	})
}

func TestRespondsWithHello(t *testing.T) {

	//New connection to Docker
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	//starting Docker container from image docker-cloud latest version
	resource, err := pool.Run("docker-cloud", "latest", []string{})
	require.NoError(t, err, "could not start container")

	//will remover container when test is complete
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	var resp *http.Response

	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("80/tcp"), "/"))
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	require.Contains(t, string(body), "Hello", "does not greet ?")
}


func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}