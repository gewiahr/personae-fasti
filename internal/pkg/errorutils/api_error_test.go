package errorutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFieldValidationErrorResponse(t *testing.T) {
	apiErr := ErrToApiError(NewFieldValidationError(map[string]string{"name": "Введите название"}))
	response := httptest.NewRecorder()
	if err := apiErr.Respond(response); err != nil {
		t.Fatalf("respond: %v", err)
	}

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	var body struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error != "validation_failed" || body.Fields["name"] != "Введите название" {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestInternalErrorSeparatesPublicMessageAndCause(t *testing.T) {
	cause := errors.New("database connection refused")
	apiErr := ErrToApiError(NewInternalError("Не удалось загрузить данные", cause))
	if apiErr.Message != "Не удалось загрузить данные" || apiErr.Code != "internal_error" {
		t.Fatalf("unexpected public error: %+v", apiErr)
	}
	if !errors.Is(apiErr.Internal(), cause) {
		t.Fatalf("internal cause was lost: %v", apiErr.Internal())
	}
}
