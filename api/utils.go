package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (api *APIServer) checkTGUserChatMembership(userID int64) (bool, error) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getChatMember?chat_id=%s&user_id=%d",
		api.auth.BotToken, "@dierolled", userID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var data map[string]any
	json.NewDecoder(resp.Body).Decode(&data)

	if !data["ok"].(bool) {
		return false, fmt.Errorf("API error: %s", data["description"])
	}

	result := data["result"].(map[string]any)
	status := result["status"].(string)

	return status == "creator" || status == "administrator" ||
		status == "member" || status == "restricted", nil
}
