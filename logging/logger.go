// Package logging provides structured logging functionality using slog with colored output.
// It supports environment-based configuration and context-aware logging patterns.
// The package is heavily inspired by: https://github.com/google/exposure-notifications-server/blob/main/pkg/logging/logger.go
//
// Example usage:
//
//	// Create a logger with specific level
//	logger := logging.NewLogger("debug", true)
//	logger.Info("Hello, world!")
//
//	// Create logger from environment variables
//	envLogger := logging.NewLoggerFromEnv()
//	envLogger.Warn("This is a warning")
//
//	// Use context-aware logging
//	ctx := logging.WithLogger(context.Background(), logger)
//	ctxLogger := logging.FromContext(ctx)
//	ctxLogger.Error("Error from context")
//
//	// Use default logger
//	defaultLogger := logging.DefaultLogger()
//	defaultLogger.Info("Using default logger")
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// NewLogger creates a new structured logger with the specified log level and development mode.
// The level parameter accepts "debug", "info", "warn", or "error" (defaults to "info" if invalid).
// The development parameter determines if the logger should use development-friendly output formatting.
// Returns a configured slog.Logger instance with colored output using the tint handler.
func NewLogger(level string, development bool) *slog.Logger {
	w := os.Stderr
	options := &tint.Options{
		TimeFormat: time.RFC3339,
	}

	switch level {
	case "debug":
		options.Level = slog.LevelDebug
	case "info":
		options.Level = slog.LevelInfo
	case "warn":
		options.Level = slog.LevelWarn
	case "error":
		options.Level = slog.LevelError
	default:
		options.Level = slog.LevelInfo
	}

	logger := slog.New(tint.NewHandler(w, options))

	return logger
}

// NewLoggerFromEnv creates a new logger from environment variables.
// It reads LOG_LEVEL to determine the logging level and APP_ENV to determine development mode.
// If APP_ENV is set to "development", development mode is enabled for better formatting.
// Automatically loads environment variables from .env file if present.
func NewLoggerFromEnv() *slog.Logger {
	_ = godotenv.Load()

	level := os.Getenv("LOG_LEVEL")
	development := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) == "development"

	return NewLogger(level, development)
}

// DefaultLogger returns the default logger for the package.
// The logger is initialized once using environment variables and cached for subsequent calls.
// This is safe for concurrent use and ensures consistent logging configuration across the application.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
// This allows for context-aware logging throughout the application call chain.
// The logger can be retrieved later using FromContext.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context.
// If no logger exists in the context, returns the default logger for the package.
// This ensures that logging is always available without needing to check for nil.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}
