package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"

	"github.com/bilte-co/toolshed/internal/cli"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// CLI represents the main command line interface
type CLI struct {
	Verbose  bool             `short:"v" help:"Enable verbose logging"`
	Version  kong.VersionFlag `help:"Show version information"`
	AES      cli.AESCmd       `cmd:"" help:"AES encryption operations"`
	Bishop   cli.BishopCmd    `cmd:"" help:"Generate ASCII art using drunken bishop algorithm"`
	Encode   cli.EncodeCmd    `cmd:"" help:"Text encoding/decoding operations"`
	Haiku    cli.HaikuCmd     `cmd:"" help:"Haiku commands"`
	Hash     cli.HashCmd      `cmd:"" help:"Hash operations"`
	Password cli.PasswordCmd  `cmd:"" help:"Password operations"`
	Serve    cli.ServeCmd     `cmd:"" help:"Start HTTP static file server"`
	ULID     cli.ULIDCmd      `cmd:"" help:"ULID operations"`
}

func main() {
	var cliApp CLI

	parser, err := kong.New(&cliApp,
		kong.Name("toolshed"),
		kong.Description("A CLI for Go utility operations"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version,
			"commit":  commit,
			"date":    date,
		},
	)
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.FatalIfErrorf(err)
	}

	// Configure logging
	setupLogging(cliApp.Verbose)

	// Execute the command
	cliContext := &cli.CLIContext{
		Logger: slog.Default(),
	}
	err = ctx.Run(cliContext)
	ctx.FatalIfErrorf(err)
}

// setupLogging configures structured logging with color output
func setupLogging(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	var output io.Writer = os.Stderr
	if !isTerminal(output) {
		// Plain text for non-terminals
		handler := slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: level,
		})
		slog.SetDefault(slog.New(handler))
		return
	}

	// Colored output for terminals
	handler := tint.NewHandler(output, &tint.Options{
		Level:      level,
		TimeFormat: "15:04:05",
	})
	slog.SetDefault(slog.New(handler))
}

// isTerminal checks if the writer is a terminal
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}
