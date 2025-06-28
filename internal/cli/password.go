package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bilte-co/toolshed/password"
)

// PasswordCmd represents the password command group
type PasswordCmd struct {
	Check PasswordCheckCmd `cmd:"" help:"Check password strength"`
}

// PasswordCheckCmd checks password strength
type PasswordCheckCmd struct {
	Text    string  `arg:"" optional:"" help:"Password to check (use '-' for stdin)"`
	Entropy float64 `long:"entropy" help:"Custom minimum entropy requirement (default: 60.0)"`
}

func (cmd *PasswordCheckCmd) Run(ctx *CLIContext) error {
	var passwordText string
	var err error

	// Handle input source
	if cmd.Text == "" || cmd.Text == "-" {
		ctx.Logger.Debug("Reading password from stdin")
		passwordText, err = cmd.readPasswordFromStdin()
		if err != nil {
			ctx.Logger.Error("Failed to read password from stdin", "error", err)
			return fmt.Errorf("failed to read password from stdin: %w", err)
		}
	} else {
		passwordText = cmd.Text
	}

	// Validate input early
	if passwordText == "" {
		ctx.Logger.Error("Password cannot be empty")
		return fmt.Errorf("password cannot be empty")
	}

	ctx.Logger.Debug("Checking password strength", "length", len(passwordText), "entropy", cmd.Entropy)

	// Use custom entropy if provided, otherwise use the password package's default
	var valid bool
	var message string

	if cmd.Entropy > 0 {
		valid, message = password.CheckEntropy(passwordText, cmd.Entropy)
		ctx.Logger.Debug("Using custom entropy", "minimum", cmd.Entropy)
	} else {
		valid, message = password.Check(passwordText)
		ctx.Logger.Debug("Using default entropy", "minimum", password.DefaultEntropy)
	}

	// Output results
	if valid {
		ctx.Logger.Info("Password validation successful")
		fmt.Println("✓ Password strength is sufficient")

		// If custom entropy was used, show the requirement
		if cmd.Entropy > 0 {
			fmt.Printf("  Meets minimum entropy requirement: %.1f\n", cmd.Entropy)
		} else {
			fmt.Printf("  Meets minimum entropy requirement: %.1f\n", password.DefaultEntropy)
		}
		return nil
	}

	// Password failed validation
	ctx.Logger.Warn("Password validation failed", "reason", message)
	fmt.Println("✗ Password strength is insufficient")
	fmt.Printf("  Reason: %s\n", message)

	// Show entropy requirement that was used
	if cmd.Entropy > 0 {
		fmt.Printf("  Required entropy: %.1f\n", cmd.Entropy)
	} else {
		fmt.Printf("  Required entropy: %.1f\n", password.DefaultEntropy)
	}

	// Exit with non-zero code on validation failure
	os.Exit(1)
	return nil
}

// readPasswordFromStdin reads a password from stdin
// It handles both piped input and terminal input
func (cmd *PasswordCheckCmd) readPasswordFromStdin() (string, error) {
	// Check if input is available from pipe
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}

	// If stdin is a pipe or file, read from it
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return cmd.readFromPipe()
	}

	// If stdin is a terminal, prompt for input
	fmt.Print("Enter password to check: ")
	return cmd.readFromTerminal()
}

// readFromPipe reads password from piped input
func (cmd *PasswordCheckCmd) readFromPipe() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from pipe: %w", err)
	}

	// Trim whitespace and newlines
	password := strings.TrimSpace(string(input))
	if password == "" {
		return "", fmt.Errorf("no password provided in piped input")
	}

	return password, nil
}

// readFromTerminal reads password from terminal input
func (cmd *PasswordCheckCmd) readFromTerminal() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read from terminal: %w", err)
		}
		return "", fmt.Errorf("no input provided")
	}

	password := strings.TrimSpace(scanner.Text())
	if password == "" {
		return "", fmt.Errorf("no password provided")
	}

	return password, nil
}

// Validate validates the command arguments
func (cmd *PasswordCheckCmd) Validate() error {
	// Validate entropy value if provided
	if cmd.Entropy < 0 {
		return fmt.Errorf("entropy value must be non-negative, got: %s", strconv.FormatFloat(cmd.Entropy, 'f', 1, 64))
	}

	// Entropy values above 200 are likely unrealistic
	if cmd.Entropy > 200 {
		return fmt.Errorf("entropy value is unrealistically high, got: %s", strconv.FormatFloat(cmd.Entropy, 'f', 1, 64))
	}

	return nil
}
