package logging_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/bilte-co/toolshed/logging"
	"github.com/stretchr/testify/require"
)

func TestNewLogger_ValidLevels(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}
	
	for _, level := range validLevels {
		logger := logging.NewLogger(level, false)
		require.NotNil(t, logger)
		require.IsType(t, &slog.Logger{}, logger)
	}
}

func TestNewLogger_InvalidLevel_DefaultsToInfo(t *testing.T) {
	logger := logging.NewLogger("invalid", false)
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLogger_DevelopmentMode(t *testing.T) {
	logger := logging.NewLogger("info", true)
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLoggerFromEnv_WithoutEnv(t *testing.T) {
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLoggerFromEnv_WithLogLevel(t *testing.T) {
	oldLogLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", oldLogLevel)
	
	os.Setenv("LOG_LEVEL", "debug")
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
}

func TestNewLoggerFromEnv_WithAppEnvDevelopment(t *testing.T) {
	oldAppEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", oldAppEnv)
	
	os.Setenv("APP_ENV", "development")
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
}

func TestNewLoggerFromEnv_WithAppEnvProduction(t *testing.T) {
	oldAppEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", oldAppEnv)
	
	os.Setenv("APP_ENV", "production")
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
}

func TestDefaultLogger_IsSingleton(t *testing.T) {
	logger1 := logging.DefaultLogger()
	logger2 := logging.DefaultLogger()
	require.Same(t, logger1, logger2)
}

func TestDefaultLogger_IsNotNil(t *testing.T) {
	logger := logging.DefaultLogger()
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestWithLogger_StoresLoggerInContext(t *testing.T) {
	logger := logging.NewLogger("info", false)
	ctx := context.Background()
	
	ctxWithLogger := logging.WithLogger(ctx, logger)
	require.NotNil(t, ctxWithLogger)
	
	retrievedLogger := logging.FromContext(ctxWithLogger)
	require.Same(t, logger, retrievedLogger)
}

func TestFromContext_WithNoLogger_ReturnsDefault(t *testing.T) {
	ctx := context.Background()
	logger := logging.FromContext(ctx)
	
	require.NotNil(t, logger)
	require.Same(t, logging.DefaultLogger(), logger)
}

func TestFromContext_WithLogger_ReturnsStoredLogger(t *testing.T) {
	customLogger := logging.NewLogger("debug", true)
	ctx := logging.WithLogger(context.Background(), customLogger)
	
	retrievedLogger := logging.FromContext(ctx)
	require.Same(t, customLogger, retrievedLogger)
}

func TestContextualLogging_RoundTrip(t *testing.T) {
	logger := logging.NewLogger("warn", false)
	ctx := context.Background()
	
	// Store logger in context
	ctxWithLogger := logging.WithLogger(ctx, logger)
	
	// Retrieve logger from context
	retrievedLogger := logging.FromContext(ctxWithLogger)
	
	// Should be the same logger
	require.Same(t, logger, retrievedLogger)
}

func TestLoggingFunctionality_CanLog(t *testing.T) {
	logger := logging.NewLogger("debug", false)
	
	// These should not panic
	require.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}

func TestLoggingWithContext_CanLog(t *testing.T) {
	logger := logging.NewLogger("info", true)
	ctx := logging.WithLogger(context.Background(), logger)
	ctxLogger := logging.FromContext(ctx)
	
	// These should not panic
	require.NotPanics(t, func() {
		ctxLogger.Info("contextual info message")
		ctxLogger.Error("contextual error message")
	})
}

func TestEnvironmentVariableHandling_CaseInsensitive(t *testing.T) {
	oldAppEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", oldAppEnv)
	
	testCases := []string{"DEVELOPMENT", "Development", "development", "  development  "}
	
	for _, testCase := range testCases {
		os.Setenv("APP_ENV", testCase)
		logger := logging.NewLoggerFromEnv()
		require.NotNil(t, logger)
	}
}
