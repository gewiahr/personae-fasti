package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func ReadBody(r *http.Request) []byte {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func ReadJsonBody(r *http.Request, v any) error {
	bodyBytes := ReadBody(r)
	return json.Unmarshal(bodyBytes, v)
}

func getPathValueInt(r *http.Request, param string) int {
	wrongValue := -1

	value := r.PathValue(param)
	if len(value) != 0 {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			return valueInt
		}
	}

	return wrongValue
}
