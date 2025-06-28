package cli

import (
	"log/slog"
)

// CLIContext provides shared context for CLI commands
type CLIContext struct {
	Logger *slog.Logger
}
