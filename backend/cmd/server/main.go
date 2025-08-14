package main

import (
	"log"
	"os"

	"github.com/noueii/nocs-log-generator/backend/pkg/api"
)

func main() {
	// Initialize router with all routes and middleware
	router := api.SetupRouter()

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("CS2 Log Generator API starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /health - Health check")
	log.Printf("  GET  /ready - Readiness check")
	log.Printf("  POST /api/v1/generate - Generate match logs")
	log.Printf("  GET  /api/v1/config/templates - Get configuration templates")
	log.Printf("  GET  /api/v1/config/maps - Get available maps")
	log.Printf("  GET  /api/v1/sample/request - Get sample request data")
	log.Printf("  GET  /api/v1/ping - API ping")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}