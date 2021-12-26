package utils

import (
	"log"

	"go.uber.org/zap"
)

// SetUpLogger uses configs to generate a logger based on
func SetUpLogger() (*zap.Logger, error) {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return logger, nil
}
