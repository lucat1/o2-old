package shared

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

// InitializeLogger initializes the logger
func InitializeLogger() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Could not initialize the logger", err)
	}
	logger = l
}

// GetLogger returns the current logger
func GetLogger() *zap.Logger {
	return logger
}
