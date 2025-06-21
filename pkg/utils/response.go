package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"message": message,
	}
	if data != nil {
		response["data"] = data
	}
	json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, message, nil)
}
