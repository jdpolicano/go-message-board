package main

import (
	"fmt"
	"net/http"

	"github.com/jdpolicano/go-message-board/internal/routes"
)

func main() {
	if err := setupServer(); err != nil {
		fmt.Println("error starting server", err)
	}
}

func setupServer() error {
	http.HandleFunc("GET /health", routes.HealthHandler)
	fmt.Println("begin listening on port 8080...")
	return http.ListenAndServe(":8080", nil)
}
