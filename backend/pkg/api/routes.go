package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noueii/nocs-log-generator/backend/pkg/websocket"
)

// SetupRouter creates and configures the main router
func SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode) // Change to gin.DebugMode for development
	
	// Create router with default middleware
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(RequestLoggingMiddleware())
	
	// Health check endpoints (not versioned)
	router.GET("/health", HealthCheckHandler)
	router.GET("/ready", ReadinessHandler)
	
	// Create WebSocket manager
	wsManager := websocket.NewManager()
	
	// Create API handler with WebSocket manager
	handler := NewHandler()
	handler.SetWebSocketManager(wsManager)
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		handler.RegisterRoutes(v1)
		
		// WebSocket endpoint
		v1.GET("/ws", wsManager.HandleWebSocketUpgrade)
	}
	
	return router
}

// HealthCheckHandler returns basic health status
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"service":   "cs2-log-generator",
		"version":   "0.1.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessHandler returns service readiness status
func ReadinessHandler(c *gin.Context) {
	// TODO: Add actual readiness checks (database, dependencies, etc.)
	c.JSON(200, gin.H{
		"status": "ready",
		"checks": gin.H{
			"api":       "ok",
			"generator": "ok", // TODO: Check generator service
			"parser":    "ok", // TODO: Check parser service
		},
	})
}

// CORSMiddleware adds CORS headers for frontend development
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLoggingMiddleware logs incoming requests
func RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Custom log format
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output:    gin.DefaultWriter,
		SkipPaths: []string{"/health", "/ready"}, // Skip health check logs to reduce noise
	})
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v", err)
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		
		c.Next()
		
		// Handle any errors that were set during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Request error: %v", err)
			
			// Don't override response if already set
			if c.Writer.Written() {
				return
			}
			
			c.JSON(500, gin.H{
				"error": "Request processing failed",
			})
		}
	}
}

// AuthMiddleware provides basic authentication (placeholder)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement authentication logic
		// For MVP, skip authentication
		c.Next()
	}
}

// RateLimitMiddleware provides rate limiting (placeholder)
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting
		// For MVP, skip rate limiting
		c.Next()
	}
}