package main

import (
	"log"
	"net/http"

	"bot_logic/handlers"
	"bot_logic/storage"
	"bot_logic/update"
)

func main() {
	if err := storage.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	handlers.RegisterRoutes()
	update.RefreshUsersTokens()

	addr := ":8083"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
