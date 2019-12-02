package utils

import "go.uber.org/zap"

// Logger creates a logger that writes logs to standard output as JSON
var Logger, _ = zap.NewProduction()
