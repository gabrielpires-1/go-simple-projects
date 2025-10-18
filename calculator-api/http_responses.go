package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type StandardResponse struct {
	Result int `json:"result"`
}

type ErrorResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func writeJSONResponseSuccess(w http.ResponseWriter, result int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	response := StandardResponse{Result: result}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeJSONResquestFailed(w http.ResponseWriter, err ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}
