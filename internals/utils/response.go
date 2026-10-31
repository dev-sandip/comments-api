package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Respond(
	w http.ResponseWriter,
	status int,
	response APIResponse,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}

func Success(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	Respond(w, status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(
	w http.ResponseWriter,
	status int,
	message string,
) {
	Respond(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}
