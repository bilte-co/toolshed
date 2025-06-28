package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bilte-co/toolshed/ulid"
)

// ULIDCmd represents the ULID command group
type ULIDCmd struct {
	Create    ULIDCreateCmd    `cmd:"" help:"Create a new ULID"`
	Timestamp ULIDTimestampCmd `cmd:"" help:"Extract timestamp from ULID"`
}

// ULIDCreateCmd creates a new ULID
type ULIDCreateCmd struct {
	Prefix    string `short:"p" help:"Optional prefix for the ULID"`
	Timestamp string `short:"t" help:"Custom timestamp (RFC3339 format, defaults to now)"`
}

// Run executes the ULID create command
func (cmd *ULIDCreateCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Creating ULID", "prefix", cmd.Prefix, "timestamp", cmd.Timestamp)

	// Validate and parse timestamp
	var timestamp time.Time
	var err error

	if cmd.Timestamp != "" {
		// Parse custom timestamp
		timestamp, err = time.Parse(time.RFC3339, cmd.Timestamp)
		if err != nil {
			ctx.Logger.Error("Invalid timestamp format", "timestamp", cmd.Timestamp, "error", err)
			return fmt.Errorf("invalid timestamp format (expected RFC3339): %w", err)
		}
	} else {
		// Use current time
		timestamp = time.Now()
	}

	// Validate prefix (basic sanitization)
	if cmd.Prefix != "" {
		if strings.ContainsAny(cmd.Prefix, " \t\n\r_") {
			ctx.Logger.Error("Invalid prefix", "prefix", cmd.Prefix)
			return fmt.Errorf("prefix cannot contain whitespace or underscore characters")
		}
		if len(cmd.Prefix) > 32 {
			ctx.Logger.Error("Prefix too long", "prefix", cmd.Prefix, "length", len(cmd.Prefix))
			return fmt.Errorf("prefix cannot exceed 32 characters")
		}
	}

	// Create ULID
	id, err := ulid.CreateULID(cmd.Prefix, timestamp)
	if err != nil {
		ctx.Logger.Error("Failed to create ULID", "error", err)
		return fmt.Errorf("failed to create ULID: %w", err)
	}

	fmt.Println(id)
	ctx.Logger.Info("ULID created successfully", "ulid", id, "prefix", cmd.Prefix)
	return nil
}

// ULIDTimestampCmd extracts timestamp from ULID
type ULIDTimestampCmd struct {
	Text   string `arg:"" help:"ULID string to decode"`
	Format string `short:"f" default:"rfc3339" help:"Output format (rfc3339, unix, unixmilli)"`
}

// Run executes the ULID timestamp command
func (cmd *ULIDTimestampCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Extracting timestamp from ULID", "ulid", cmd.Text, "format", cmd.Format)

	// Read from stdin if text is "-"
	var input string
	if cmd.Text == "-" {
		data, err := readStdin()
		if err != nil {
			ctx.Logger.Error("Failed to read from stdin", "error", err)
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		input = strings.TrimSpace(string(data))
	} else {
		input = cmd.Text
	}

	// Validate input
	if input == "" {
		ctx.Logger.Error("Empty ULID input")
		return fmt.Errorf("ULID cannot be empty")
	}

	// Extract timestamp
	timestamp, err := ulid.Timestamp(input)
	if err != nil {
		ctx.Logger.Error("Failed to extract timestamp", "ulid", input, "error", err)
		return fmt.Errorf("failed to extract timestamp: %w", err)
	}

	// Format output
	var output string
	switch strings.ToLower(cmd.Format) {
	case "rfc3339":
		output = timestamp.Format(time.RFC3339)
	case "unix":
		output = fmt.Sprintf("%d", timestamp.Unix())
	case "unixmilli":
		output = fmt.Sprintf("%d", timestamp.UnixMilli())
	default:
		ctx.Logger.Error("Invalid format", "format", cmd.Format)
		return fmt.Errorf("invalid format '%s' (supported: rfc3339, unix, unixmilli)", cmd.Format)
	}

	fmt.Println(output)
	ctx.Logger.Info("Timestamp extracted successfully", "ulid", input, "timestamp", output)
	return nil
}

// readStdin reads data from standard input
func readStdin() ([]byte, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	// Check if there's data available
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data available from stdin")
	}

	// Read all data from stdin
	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
	}

	return data, nil
}
