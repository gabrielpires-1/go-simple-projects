package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type StandardResponse struct {
	Result int `json:"result"`
}

func writeJSONResponse(w http.ResponseWriter, result int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	response := StandardResponse{Result: result}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
