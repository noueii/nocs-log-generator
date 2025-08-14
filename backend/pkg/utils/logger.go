package utils

import (
	"log"
	"os"
)

// Logger provides structured logging functionality
type Logger struct {
	*log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[CS2-LOG-GEN] ", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs info level messages
func (l *Logger) Info(msg string) {
	l.Printf("INFO: %s", msg)
}

// Error logs error level messages
func (l *Logger) Error(msg string) {
	l.Printf("ERROR: %s", msg)
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string) {
	l.Printf("DEBUG: %s", msg)
}