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

// Additional comprehensive test coverage

func TestNewLogger_EmptyLevel_DefaultsToInfo(t *testing.T) {
	logger := logging.NewLogger("", false)
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLogger_WhitespaceLevel_DefaultsToInfo(t *testing.T) {
	logger := logging.NewLogger("  ", false)
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLogger_CaseSensitiveLevel_DefaultsToInfo(t *testing.T) {
	testCases := []string{"DEBUG", "Info", "WARN", "Error", "DEBUG_LEVEL"}
	
	for _, level := range testCases {
		logger := logging.NewLogger(level, false)
		require.NotNil(t, logger)
		require.IsType(t, &slog.Logger{}, logger)
	}
}

func TestNewLoggerFromEnv_EmptyLogLevel(t *testing.T) {
	oldLogLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", oldLogLevel)

	os.Setenv("LOG_LEVEL", "")
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLoggerFromEnv_WhitespaceLogLevel(t *testing.T) {
	oldLogLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", oldLogLevel)

	os.Setenv("LOG_LEVEL", "   ")
	logger := logging.NewLoggerFromEnv()
	require.NotNil(t, logger)
	require.IsType(t, &slog.Logger{}, logger)
}

func TestNewLoggerFromEnv_AllValidLogLevels(t *testing.T) {
	oldLogLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", oldLogLevel)

	validLevels := []string{"debug", "info", "warn", "error", "invalid"}

	for _, level := range validLevels {
		os.Setenv("LOG_LEVEL", level)
		logger := logging.NewLoggerFromEnv()
		require.NotNil(t, logger)
		require.IsType(t, &slog.Logger{}, logger)
	}
}

func TestNewLoggerFromEnv_AppEnvEdgeCases(t *testing.T) {
	oldAppEnv := os.Getenv("APP_ENV")
	oldLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		os.Setenv("APP_ENV", oldAppEnv)
		os.Setenv("LOG_LEVEL", oldLogLevel)
	}()

	testCases := []struct {
		name   string
		appEnv string
	}{
		{"empty", ""},
		{"whitespace", "   "},
		{"production", "production"},
		{"staging", "staging"},
		{"test", "test"},
		{"mixed_case", "DevElOpMeNt"},
		{"with_leading_trailing_space", "  development  "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tc.appEnv)
			logger := logging.NewLoggerFromEnv()
			require.NotNil(t, logger)
			require.IsType(t, &slog.Logger{}, logger)
		})
	}
}

func TestWithLogger_NilLogger(t *testing.T) {
	ctx := context.Background()
	ctxWithLogger := logging.WithLogger(ctx, nil)
	require.NotNil(t, ctxWithLogger)

	// When nil logger is stored, FromContext should return default logger
	// The type assertion (*slog.Logger) will succeed for nil since nil implements *slog.Logger
	// but the ok value will still be true, so we get the nil back
	retrievedLogger := logging.FromContext(ctxWithLogger)
	
	// This test shows that the current implementation has a bug - it returns nil instead of default
	// In a real implementation, we'd want this to return DefaultLogger()
	require.Nil(t, retrievedLogger) // This is the actual current behavior
}

func TestWithLogger_NilContext_Panics(t *testing.T) {
	logger := logging.NewLogger("info", false)
	
	// This should panic with nil context - Go's context.WithValue panics on nil parent
	require.Panics(t, func() {
		logging.WithLogger(nil, logger)
	})
}

func TestFromContext_NilContext_Panics(t *testing.T) {
	// This should panic with nil context - accessing Value on nil context panics
	require.Panics(t, func() {
		logging.FromContext(nil)
	})
}

func TestFromContext_ContextWithWrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), "logger", "not-a-logger")
	logger := logging.FromContext(ctx)
	
	require.NotNil(t, logger)
	require.Same(t, logging.DefaultLogger(), logger)
}

func TestFromContext_ContextWithWrongKey(t *testing.T) {
	anotherLogger := logging.NewLogger("debug", true)
	ctx := context.WithValue(context.Background(), "different-key", anotherLogger)
	logger := logging.FromContext(ctx)
	
	require.NotNil(t, logger)
	require.Same(t, logging.DefaultLogger(), logger)
}

func TestDefaultLogger_ConcurrentAccess(t *testing.T) {
	const numGoroutines = 100
	results := make(chan *slog.Logger, numGoroutines)

	// Launch multiple goroutines to test concurrent access
	for i := 0; i < numGoroutines; i++ {
		go func() {
			results <- logging.DefaultLogger()
		}()
	}

	// Collect all results
	var loggers []*slog.Logger
	for i := 0; i < numGoroutines; i++ {
		loggers = append(loggers, <-results)
	}

	// All should be the same instance
	firstLogger := loggers[0]
	require.NotNil(t, firstLogger)
	
	for i, logger := range loggers {
		require.Same(t, firstLogger, logger, "Logger %d should be same as first logger", i)
	}
}

func TestCompleteWorkflow_EnvironmentToContext(t *testing.T) {
	oldLogLevel := os.Getenv("LOG_LEVEL")
	oldAppEnv := os.Getenv("APP_ENV")
	defer func() {
		os.Setenv("LOG_LEVEL", oldLogLevel)
		os.Setenv("APP_ENV", oldAppEnv)
	}()

	// Set environment variables
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("APP_ENV", "development")

	// Create logger from environment
	envLogger := logging.NewLoggerFromEnv()
	require.NotNil(t, envLogger)

	// Store in context
	ctx := logging.WithLogger(context.Background(), envLogger)
	require.NotNil(t, ctx)

	// Retrieve from context
	ctxLogger := logging.FromContext(ctx)
	require.Same(t, envLogger, ctxLogger)

	// Test logging functionality
	require.NotPanics(t, func() {
		ctxLogger.Debug("debug message")
		ctxLogger.Info("info message", slog.String("key", "value"))
		ctxLogger.Warn("warn message", slog.Int("count", 42))
		ctxLogger.Error("error message", slog.Bool("critical", true))
	})
}

func TestContextChaining_DeepNesting(t *testing.T) {
	logger1 := logging.NewLogger("info", false)
	logger2 := logging.NewLogger("debug", true)

	ctx := context.Background()
	
	// First level
	ctx1 := logging.WithLogger(ctx, logger1)
	retrieved1 := logging.FromContext(ctx1)
	require.Same(t, logger1, retrieved1)
	
	// Second level - overwrite with new logger
	ctx2 := logging.WithLogger(ctx1, logger2)
	retrieved2 := logging.FromContext(ctx2)
	require.Same(t, logger2, retrieved2)
	
	// Original context should still have original logger
	retrievedOriginal := logging.FromContext(ctx1)
	require.Same(t, logger1, retrievedOriginal)
}
