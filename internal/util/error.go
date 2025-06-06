package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if e := json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	}); e != nil {
		log.Fatal(e)
	}
}
