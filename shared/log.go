package shared

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

// InitializeLogger initializes the logger
func InitializeLogger() {
	l, err := zap.NewProduction()
	if err != nil {
		panic("Could not initialize the logger: \n" + err.Error())
	}
	logger = l
}

// GetLogger returns the current logger
func GetLogger() *zap.Logger {
	return logger
}
