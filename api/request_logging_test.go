package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/service"
)

type recordingLogRepo struct {
	logs []*domain.ApiLog
	err  error
}

func (r *recordingLogRepo) Insert(_ context.Context, entry *domain.ApiLog) error {
	r.logs = append(r.logs, entry)
	return r.err
}

func TestLogRequestsCapturesAndRedactsJSON(t *testing.T) {
	repo := &recordingLogRepo{}
	logger := service.NewLogService(repo)
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setRequestLogPlayer(r.Context(), 42)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"token":"response-token","name":"visible"}`))
	}), logger)

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"loginData":"password","username":"visible"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if len(repo.logs) != 1 {
		t.Fatalf("logged requests = %d, want 1", len(repo.logs))
	}
	entry := repo.logs[0]
	if entry.PlayerID != 42 || entry.Code != http.StatusCreated || entry.URI != "/login" {
		t.Fatalf("unexpected log metadata: %+v", entry)
	}
	if strings.Contains(entry.Request, "password") || strings.Contains(entry.Response, "response-token") {
		t.Fatalf("sensitive data was logged: request=%q response=%q", entry.Request, entry.Response)
	}
	if !strings.Contains(entry.Request, "visible") || !strings.Contains(entry.Response, "visible") {
		t.Fatalf("non-sensitive data was removed: request=%q response=%q", entry.Request, entry.Response)
	}
}

func TestLogRequestsSkipsHealthChecks(t *testing.T) {
	repo := &recordingLogRepo{}
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), service.NewLogService(repo))

	for _, path := range []string{"/healthz", "/readyz"} {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}
	if len(repo.logs) != 0 {
		t.Fatalf("health logs = %d, want 0", len(repo.logs))
	}
}

func TestLogRequestsRedactsPersonalNotes(t *testing.T) {
	repo := &recordingLogRepo{}
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"text":"private response"}`))
	}), service.NewLogService(repo))

	request := httptest.NewRequest(http.MethodPut, "/notes", strings.NewReader(`{"text":"private request"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(httptest.NewRecorder(), request)

	entry := repo.logs[0]
	if entry.Request != "[redacted: personal note]" || entry.Response != "[redacted: personal note]" {
		t.Fatalf("personal note was not redacted: request=%q response=%q", entry.Request, entry.Response)
	}
}

func TestLogRequestsDoesNotChangeResponseOnLogFailure(t *testing.T) {
	repo := &recordingLogRepo{err: errors.New("database unavailable")}
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}), service.NewLogService(repo))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/record/1", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("response status = %d, want %d", response.Code, http.StatusNoContent)
	}
}
