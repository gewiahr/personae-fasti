package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"personae-fasti/internal/domain"
	"personae-fasti/internal/service"
)

const maxLoggedBodyBytes = 64 << 10

type requestLogContextKey struct{}

type requestLogState struct {
	requestID     string
	playerID      int
	gameID        *int
	errorCode     string
	internalError string
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
		state := &requestLogState{requestID: newRequestID()}
		r = r.WithContext(context.WithValue(r.Context(), requestLogContextKey{}, state))
		w.Header().Set("X-Request-ID", state.requestID)
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
		completedAt := time.Now().UTC()
		publicError := ""
		if status >= http.StatusBadRequest {
			publicError = responseBody
			if state.errorCode == "" {
				state.errorCode = "http_error"
			}
		}
		entry := &domain.ApiLog{
			PlayerID: optionalInt(state.playerID), GameID: state.gameID, RequestID: state.requestID,
			IP: optionalString(clientIP(r)), Host: optionalString(r.Host),
			URI: r.URL.Path, Method: r.Method, Request: optionalString(requestBody),
			Response: optionalString(responseBody), Code: status, Error: optionalString(publicError),
			ErrorCode: optionalString(state.errorCode), InternalError: optionalString(state.internalError),
			Time: completedAt.Sub(startedAt).Milliseconds(), Started: startedAt, Created: completedAt,
		}
		if err := logService.InsertLog(entry); err != nil {
			log.Printf("failed to insert API log for %s %s: %v", r.Method, r.URL.Path, err)
		}
	})
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func optionalInt(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

func setRequestLogPlayer(ctx context.Context, playerID, gameID int) {
	if state, ok := ctx.Value(requestLogContextKey{}).(*requestLogState); ok {
		state.playerID = playerID
		if gameID != 0 {
			state.gameID = &gameID
		}
	}
}

func setRequestLogError(ctx context.Context, code string, err error) {
	if state, ok := ctx.Value(requestLogContextKey{}).(*requestLogState); ok {
		state.errorCode = code
		if err != nil {
			state.internalError = err.Error()
		}
	}
}

func requestID(ctx context.Context) string {
	if state, ok := ctx.Value(requestLogContextKey{}).(*requestLogState); ok {
		return state.requestID
	}
	return ""
}

func newRequestID() string {
	data := make([]byte, 12)
	if _, err := rand.Read(data); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(data)
}

func captureRequestBody(r *http.Request) string {
	if r.Body == nil || r.ContentLength <= 0 {
		return ""
	}
	if r.ContentLength > maxLoggedBodyBytes {
		return "[omitted: body too large]"
	}
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "multipart/") || strings.Contains(contentType, "octet-stream") || strings.HasPrefix(contentType, "image/") {
		return "[omitted: non-JSON body]"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Sprintf("[omitted: failed to read body: %v]", err)
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	if !json.Valid(body) {
		return "[omitted: non-JSON body]"
	}
	return sanitizeLoggedBody(r.URL.Path, body)
}

func clientIP(r *http.Request) string {
	// Nginx overwrites X-Real-IP with the direct client address. Prefer it to
	// X-Forwarded-For, which may contain client-supplied values in its chain.
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first, _, ok := strings.Cut(forwarded, ","); ok {
			return strings.TrimSpace(first)
		}
		return strings.TrimSpace(forwarded)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
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
