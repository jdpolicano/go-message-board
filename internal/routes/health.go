package routes

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	TimeStamp int64 `json:"timeStamp"`
	Status    int   `json:"status"`
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		TimeStamp: time.Now().UnixMilli(),
		Status:    200,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
