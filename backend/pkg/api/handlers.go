package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler contains dependencies for API handlers
type Handler struct {
	// TODO: Add dependencies like generator, parser, etc.
}

// NewHandler creates a new API handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes sets up API routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// TODO: Register API endpoints here
	// router.POST("/generate", h.GenerateMatch)
	// router.POST("/parse", h.ParseDemo)
	
	// Placeholder endpoint
	router.GET("/ping", h.Ping)
}

// Ping is a simple ping endpoint for testing
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"api":     "cs2-log-generator",
	})
}