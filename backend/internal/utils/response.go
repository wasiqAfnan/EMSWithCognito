package utils

import (
	"encoding/json"
	"net/http"

	"awsems/internal/model"
)

func JSON(w http.ResponseWriter, statusCode int, message string, data any) {
	response := model.Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, message, nil)
}

func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	JSON(w, statusCode, message, data)
}
