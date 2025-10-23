package main

import (
	"chat-app-backend/internal/api"
	"chat-app-backend/internal/config"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env first
	cfg := config.LoadConfig()

	// Connect to DB
	config.ConnectDatabase(cfg)
	log.Println("Database connected successfully")

	// Setup Gin
	router := gin.Default()
	api.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	router.Run(":" + port)
}
