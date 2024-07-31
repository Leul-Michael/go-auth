package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	type SuccessMsg struct {
		Status string      `json:"status"`
		Code   int         `json:"code"`
		Data   interface{} `json:"data"`
	}

	successRes := &SuccessMsg{
		Status: "ok",
		Code:   code,
		Data:   payload,
	}

	data, err := json.Marshal(successRes)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("failed to marasal json: %v", payload)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message interface{}) {
	response := map[string]interface{}{
		"status": "fail",
		"code":   code,
	}

	switch v := message.(type) {
	case string:
		response["error"] = v
	case []string:
		response["errors"] = v
	default:
		response["error"] = "something went wrong!"
	}

	if code > 499 {
		fmt.Println("server error detected.")
	}

	data, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("failed to marasal json: %v", response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("invalid request payload: %w", err)
	}
	return v, nil
}
