package testutil

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/bilte-co/toolshed/internal/cli"
)

// NewTestContext creates a test CLI context with a logger that discards output
func NewTestContext() *cli.CLIContext {
	return &cli.CLIContext{
		Logger: slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
			Level: slog.LevelError, // Only show errors to reduce test noise
		})),
	}
}

// Mock os.Exit for testing
var OsExit = os.Exit
