package main

import (
	"fmt"
	"net/http"

	"github.com/jdpolicano/go-message-board/internal/controller"
	"github.com/jdpolicano/go-message-board/internal/db"
	"github.com/jdpolicano/go-message-board/internal/routes"
)

func main() {
	if err := setupServer(); err != nil {
		fmt.Println("error starting server", err)
	}
}

func setupServer() error {
	db := db.NewMemDatabase()
	controller := controller.NewController(&db)
	http.HandleFunc("GET /health", routes.HealthHandler)
	http.HandleFunc("GET /chat/sessions", routes.NewListChatHandler(controller))
	fmt.Println("begin listening on port 8080...")
	return http.ListenAndServe(":8080", nil)
}
