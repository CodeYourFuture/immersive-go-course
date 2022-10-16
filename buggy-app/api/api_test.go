package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/api/model"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/auth"
	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/pashagolub/pgxmock/v2"
)

var defaultConfig Config = Config{
	Port:           8090,
	Log:            log.Default(),
	AuthServiceUrl: "auth:8080",
}

func assertJSON(actual []byte, data interface{}, t *testing.T) {
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	if !bytes.Equal(expected, actual) {
		t.Errorf("the expected json: %s is different from actual %s", expected, actual)
	}
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
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
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
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
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
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateDeny,
	})

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic nope")
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestMyNotesAuthPass(t *testing.T) {
	as := New(defaultConfig)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateAllow,
	})

	rows := mock.NewRows([]string{"id", "owner", "content"})

	mock.ExpectQuery("^SELECT (.+) FROM public.note$").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic ZXhhbXBsZTpleGFtcGxl")
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestMyNotesOneNone(t *testing.T) {
	as := New(defaultConfig)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateAllow,
	})

	id, password := "abc123", "password"
	noteId, content, created, modified := "xyz789", "Note content", time.Now(), time.Now()

	rows := mock.NewRows([]string{"id", "owner", "content", "created", "modified"}).
		AddRow(noteId, id, content, created, modified)

	mock.ExpectQuery("^SELECT (.+) FROM public.note$").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", util.BasicAuthHeaderValue(id, password))
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	data := struct {
		Notes []model.Note `json:"notes"`
	}{Notes: []model.Note{
		{Id: noteId, Owner: id, Content: content, Created: created, Modified: modified, Tags: []string{}},
	}}
	assertJSON(res.Body.Bytes(), data, t)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %s", err)
	}
}

func TestMyNotesNonOwnedNote(t *testing.T) {
	as := New(defaultConfig)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateAllow,
	})

	id, password := "abc123", "password"
	noteId, content, created, modified := "xyz789", "Note content", time.Now(), time.Now()

	rows := mock.NewRows([]string{"id", "owner", "content", "created", "modified"}).
		AddRow(noteId, id, content, created, modified).
		AddRow("pqr123", "mno456", "Non-owned note", created, modified)

	mock.ExpectQuery("^SELECT (.+) FROM public.note$").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/1/my/notes.json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", util.BasicAuthHeaderValue(id, password))
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	data := struct {
		Notes []model.Note `json:"notes"`
	}{Notes: []model.Note{
		{Id: noteId, Owner: id, Content: content, Created: created, Modified: modified, Tags: []string{}},
	}}
	assertJSON(res.Body.Bytes(), data, t)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %s", err)
	}
}

func TestMyNoteById(t *testing.T) {
	as := New(defaultConfig)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateAllow,
	})

	id, password := "abc123", "password"
	noteId, content, created, modified := "xyz789", "Note content", time.Now(), time.Now()

	rows := mock.NewRows([]string{"id", "owner", "content", "created", "modified"}).
		AddRow(noteId, id, content, created, modified)

	mock.ExpectQuery("^SELECT (.+) FROM public.note WHERE id = (.+)$").WillReturnRows(rows)

	req, err := http.NewRequest("GET", fmt.Sprintf("/1/my/note/%s.json", noteId), strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", util.BasicAuthHeaderValue(id, password))
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	data := struct {
		Note model.Note `json:"note"`
	}{Note: model.Note{Id: noteId, Owner: id, Content: content, Created: created, Modified: modified, Tags: []string{}}}
	assertJSON(res.Body.Bytes(), data, t)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %s", err)
	}
}

func TestMyNoteByIdWithTags(t *testing.T) {
	as := New(defaultConfig)
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	as.pool = mock
	as.authClient = auth.NewMockClient(&auth.VerifyResult{
		State: auth.StateAllow,
	})

	id, password := "abc123", "password"
	noteId, content, created, modified := "xyz789", "Note content #tag1", time.Now(), time.Now()

	rows := mock.NewRows([]string{"id", "owner", "content", "created", "modified"}).
		AddRow(noteId, id, content, created, modified)

	mock.ExpectQuery("^SELECT (.+) FROM public.note WHERE id = (.+)$").WillReturnRows(rows)

	req, err := http.NewRequest("GET", fmt.Sprintf("/1/my/note/%s.json", noteId), strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", util.BasicAuthHeaderValue(id, password))
	res := httptest.NewRecorder()
	handler := as.Handler()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	data := struct {
		Note model.Note `json:"note"`
	}{Note: model.Note{Id: noteId, Owner: id, Content: content, Created: created, Modified: modified, Tags: []string{"tag1"}}}
	assertJSON(res.Body.Bytes(), data, t)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %s", err)
	}
}
