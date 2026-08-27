package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"personae-fasti/internal/domain"
	e "personae-fasti/internal/pkg/errorutils"
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

func (r *recordingLogRepo) Prune(_ context.Context, _, _ *time.Time) error { return r.err }

func TestLogRequestsCapturesAndRedactsJSON(t *testing.T) {
	repo := &recordingLogRepo{}
	logger := service.NewLogService(repo)
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setRequestLogPlayer(r.Context(), 42, 7)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"token":"response-token","name":"visible"}`))
	}), logger)

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"loginData":"password","username":"visible"}`))
	request.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	request.Header.Set("X-Real-IP", "203.0.113.8")
	request.Header.Set("X-Forwarded-For", "spoofed, 203.0.113.8")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if len(repo.logs) != 1 {
		t.Fatalf("logged requests = %d, want 1", len(repo.logs))
	}
	entry := repo.logs[0]
	if entry.PlayerID == nil || *entry.PlayerID != 42 || entry.Code != http.StatusCreated || entry.URI != "/login" {
		t.Fatalf("unexpected log metadata: %+v", entry)
	}
	if entry.GameID == nil || *entry.GameID != 7 || entry.RequestID == "" || entry.Started.IsZero() {
		t.Fatalf("missing structured log metadata: %+v", entry)
	}
	if stringValue(entry.IP) != "203.0.113.8" || stringValue(entry.Host) != "example.com" {
		t.Fatalf("missing client metadata: ip=%q host=%q", stringValue(entry.IP), stringValue(entry.Host))
	}
	if strings.Contains(stringValue(entry.Request), "password") || strings.Contains(stringValue(entry.Response), "response-token") {
		t.Fatalf("sensitive data was logged: request=%q response=%q", stringValue(entry.Request), stringValue(entry.Response))
	}
	if !strings.Contains(stringValue(entry.Request), "visible") || !strings.Contains(stringValue(entry.Response), "visible") {
		t.Fatalf("non-sensitive data was removed: request=%q response=%q", stringValue(entry.Request), stringValue(entry.Response))
	}
	if entry.Error != nil || entry.ErrorCode != nil || entry.InternalError != nil {
		t.Fatalf("successful request has error values: %+v", entry)
	}
}

func TestLogRequestsUsesNullForAbsentValues(t *testing.T) {
	repo := &recordingLogRepo{}
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}), service.NewLogService(repo))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/empty", nil))

	entry := repo.logs[0]
	if entry.PlayerID != nil || entry.Request != nil || entry.Response != nil || entry.Error != nil || entry.ErrorCode != nil || entry.InternalError != nil {
		t.Fatalf("absent values were not NULL: %+v", entry)
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
	if stringValue(entry.Request) != "[redacted: personal note]" || stringValue(entry.Response) != "[redacted: personal note]" {
		t.Fatalf("personal note was not redacted: request=%q response=%q", stringValue(entry.Request), stringValue(entry.Response))
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

func TestLogRequestsSeparatesPublicAndInternalErrors(t *testing.T) {
	repo := &recordingLogRepo{}
	handler := logRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondAndTrack(w, r, e.ErrToApiError(e.NewInternalError("Не удалось сохранить запись", errors.New("database constraint failed"))))
	}), service.NewLogService(repo))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/record", nil))
	entry := repo.logs[0]
	if stringValue(entry.ErrorCode) != "internal_error" || stringValue(entry.InternalError) != "database constraint failed" {
		t.Fatalf("internal error was not captured: %+v", entry)
	}
	if strings.Contains(response.Body.String(), "database constraint") || !strings.Contains(response.Body.String(), "Не удалось сохранить запись") {
		t.Fatalf("unsafe or missing public response: %q", response.Body.String())
	}
	if response.Header().Get("X-Request-ID") == "" || !strings.Contains(response.Body.String(), entry.RequestID) {
		t.Fatalf("request ID was not correlated: header=%q body=%q", response.Header().Get("X-Request-ID"), response.Body.String())
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
