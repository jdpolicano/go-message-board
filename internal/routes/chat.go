package routes

type BoardListResonse struct {
	Boards []string `json:"boards"`
}

type BoardCreateRequest struct {
	Name string `json:"name"`
}

type Chat struct {
}

// func HealthHandler(w http.ResponseWriter, _ *http.Request) {
// 	response := HealthResponse{
// 		TimeStamp: time.Now().Unix(),
// 		Status:    200,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
