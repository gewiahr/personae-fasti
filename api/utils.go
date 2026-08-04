package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	e "personae-fasti/internal/pkg/errorutils"
)

func readBody[Body any](w http.ResponseWriter, r *http.Request) (Body, *e.ApiError) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		var zero Body
		var err = e.NewApiError(http.StatusRequestEntityTooLarge, "request body too large")
		return zero, &err
	}
	defer r.Body.Close()

	var body Body
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			var zero Body
			var err = e.NewApiError(http.StatusBadRequest, "invalid JSON")
			return zero, &err
		}
	}
	return body, nil
}

func extractBearer(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return parts[len(parts)-1]
}
