package infrastructure

import (
	"time"

	"go.uber.org/zap"
)

// ZapLogger struct wraps a zap.Logger to provide custom logging functionality
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger returns an instance of ZapLogger instance configured for production logging
// It sets the output path to stdout
func NewZapLogger() (*ZapLogger, error) {
	// configuration for logging to console
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	// build the logger
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: l}, nil
}

// Info logs an informational message with the provided fields
func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	// log the message at the Info level
	z.logger.Info(msg, fields...)
}

// Sync flushes any buffered log entries
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

// String creates a string field in zap
func (z *ZapLogger) String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Time creates a time field in zap
func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}
