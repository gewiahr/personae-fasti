package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"personae-fasti/internal/service"
)

const maxLoggedBodyBytes = 64 << 10

type requestLogContextKey struct{}

type requestLogState struct {
	playerID int
}

type logResponseWriter struct {
	http.ResponseWriter
	status    int
	body      bytes.Buffer
	truncated bool
}

func (w *logResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *logResponseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	remaining := maxLoggedBodyBytes - w.body.Len()
	if remaining > 0 {
		captured := len(data)
		if captured > remaining {
			captured = remaining
			w.truncated = true
		}
		_, _ = w.body.Write(data[:captured])
	} else if len(data) > 0 {
		w.truncated = true
	}
	return w.ResponseWriter.Write(data)
}

func logRequests(next http.Handler, logService *service.LogService) http.Handler {
	if logService == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}

		startedAt := time.Now().UTC()
		state := &requestLogState{}
		r = r.WithContext(context.WithValue(r.Context(), requestLogContextKey{}, state))
		requestBody := captureRequestBody(r)
		responseWriter := &logResponseWriter{ResponseWriter: w}

		next.ServeHTTP(responseWriter, r)
		status := responseWriter.status
		if status == 0 {
			status = http.StatusOK
		}
		responseBody := sanitizeLoggedBody(r.URL.Path, responseWriter.body.Bytes())
		if responseWriter.truncated {
			responseBody += " [truncated]"
		}
		if err := logService.InsertLog(state.playerID, r, requestBody, responseBody, status, startedAt); err != nil {
			log.Printf("failed to insert API log for %s %s: %v", r.Method, r.URL.Path, err)
		}
	})
}

func setRequestLogPlayer(ctx context.Context, playerID int) {
	if state, ok := ctx.Value(requestLogContextKey{}).(*requestLogState); ok {
		state.playerID = playerID
	}
}

func captureRequestBody(r *http.Request) string {
	if r.Body == nil || r.ContentLength <= 0 {
		return ""
	}
	if r.ContentLength > maxLoggedBodyBytes {
		return "[omitted: body too large]"
	}
	if !strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
		return "[omitted: non-JSON body]"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Sprintf("[omitted: failed to read body: %v]", err)
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	return sanitizeLoggedBody(r.URL.Path, body)
}

func sanitizeLoggedBody(path string, body []byte) string {
	if len(body) == 0 {
		return ""
	}
	if path == "/notes" {
		return "[redacted: personal note]"
	}

	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}
	redactJSON(value)
	sanitized, err := json.Marshal(value)
	if err != nil {
		return "[omitted: failed to sanitize JSON]"
	}
	return string(sanitized)
}

func redactJSON(value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if isSensitiveLogField(key) {
				typed[key] = "[redacted]"
				continue
			}
			redactJSON(child)
		}
	case []any:
		for _, child := range typed {
			redactJSON(child)
		}
	}
}

func isSensitiveLogField(field string) bool {
	switch strings.ToLower(field) {
	case "authorization", "password", "passwordhash", "password_hash", "logindata", "token", "tokenhash", "token_hash", "accesskey", "secretkey":
		return true
	default:
		return false
	}
}
