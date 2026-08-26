package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"personae-fasti/configs"
)

func TestHealthEndpoints(t *testing.T) {
	server := ConfigServer(&configs.Main{App: &configs.APIConfig{}}, nil, nil)

	health := httptest.NewRecorder()
	server.Server.Handler.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", health.Code, http.StatusOK)
	}

	ready := httptest.NewRecorder()
	server.Server.Handler.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if ready.Code != http.StatusOK {
		t.Fatalf("ready status = %d, want %d", ready.Code, http.StatusOK)
	}
}

func TestReadinessFailure(t *testing.T) {
	server := ConfigServer(&configs.Main{App: &configs.APIConfig{}}, nil, func(_ context.Context) error {
		return errors.New("database unavailable")
	})

	response := httptest.NewRecorder()
	server.Server.Handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("ready status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
}

func TestProductionCORSAndUploadTimeout(t *testing.T) {
	server := ConfigServer(&configs.Main{App: &configs.APIConfig{
		Environment:        "production",
		AllowedOrigins:     []string{"https://app.storyshard.ru"},
		AllowCredentials:   true,
		ReadTimeout:        10,
		WriteTimeout:       10,
		ImageUploadTimeout: 90,
	}}, nil, nil)

	if server.Server.ReadTimeout != 90*time.Second || server.Server.WriteTimeout != 90*time.Second {
		t.Fatalf("server timeouts = %s/%s, want at least upload timeout", server.Server.ReadTimeout, server.Server.WriteTimeout)
	}

	request := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	request.Header.Set("Origin", "https://app.storyshard.ru")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	response := httptest.NewRecorder()
	server.Server.Handler.ServeHTTP(response, request)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://app.storyshard.ru" {
		t.Fatalf("allow origin = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("allow credentials = %q", got)
	}
}

func TestWildcardCORSDisablesCredentials(t *testing.T) {
	server := ConfigServer(&configs.Main{App: &configs.APIConfig{
		Debug:            true,
		AllowCredentials: true,
	}}, nil, nil)

	request := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	response := httptest.NewRecorder()
	server.Server.Handler.ServeHTTP(response, request)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("allow origin = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("allow credentials must be empty for wildcard origin, got %q", got)
	}
}
