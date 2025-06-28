package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/bilte-co/toolshed/aes"
)

// AESCmd represents the AES command group
type AESCmd struct {
	GenerateKey GenerateKeyCmd `cmd:"generate-key" help:"Generate a new AES key"`
	Encrypt     EncryptCmd     `cmd:"" help:"Encrypt a file using AES-GCM"`
	Decrypt     DecryptCmd     `cmd:"" help:"Decrypt a file using AES-GCM"`
}

// GenerateKeyCmd generates a new AES key
type GenerateKeyCmd struct {
	Entropy int `short:"e" default:"256" help:"Key entropy in bits (128, 192, or 256)"`
}

func (cmd *GenerateKeyCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating AES key", "entropy", cmd.Entropy)

	// Validate entropy
	if cmd.Entropy != 128 && cmd.Entropy != 192 && cmd.Entropy != 256 {
		err := fmt.Errorf("invalid entropy: %d (must be 128, 192, or 256)", cmd.Entropy)
		ctx.Logger.Error("Invalid entropy specified", "entropy", cmd.Entropy, "error", err)
		return err
	}

	key, err := aes.GenerateAESKey(cmd.Entropy)
	if err != nil {
		ctx.Logger.Error("Failed to generate AES key", "error", err)
		return err
	}

	fmt.Println(key)
	ctx.Logger.Info("AES key generated successfully", "entropy", cmd.Entropy)
	return nil
}

// EncryptCmd encrypts a file using AES-GCM
type EncryptCmd struct {
	File   string `arg:"" help:"File to encrypt (use '-' for stdin)"`
	Key    string `short:"k" help:"Base64-encoded AES key (if not provided, reads from AES_KEY env var)"`
	Output string `short:"o" help:"Output file (if not provided, prints to stdout)"`
}

func (cmd *EncryptCmd) Run(ctx *CLIContext) error {
	// Get the key
	key, err := cmd.getKey()
	if err != nil {
		ctx.Logger.Error("Failed to get encryption key", "error", err)
		return err
	}

	ctx.Logger.Debug("Encrypting file", "file", cmd.File, "output", cmd.Output)

	// Read input
	var input io.Reader
	var inputName string

	if cmd.File == "-" {
		input = os.Stdin
		inputName = "stdin"
		ctx.Logger.Debug("Reading from stdin")
	} else {
		// Sanitize and validate file path
		cleanPath := filepath.Clean(cmd.File)
		if !strings.HasPrefix(cleanPath, "/") && !strings.HasPrefix(cleanPath, "./") && !strings.HasPrefix(cleanPath, "../") {
			cleanPath = "./" + cleanPath
		}

		file, err := os.Open(cleanPath)
		if err != nil {
			ctx.Logger.Error("Failed to open input file", "path", cleanPath, "error", err)
			return fmt.Errorf("failed to open input file %s: %w", cleanPath, err)
		}
		defer file.Close()

		input = file
		inputName = cleanPath
	}

	// Show spinner for file operations (only if not reading from stdin and not outputting to stdout)
	var s *spinner.Spinner
	if cmd.File != "-" && cmd.Output != "" {
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Encrypting file..."
		s.Start()
		defer s.Stop()
	}

	// Read all data - for production, consider streaming encryption for large files
	data, err := io.ReadAll(input)
	if err != nil {
		ctx.Logger.Error("Failed to read input", "source", inputName, "error", err)
		return fmt.Errorf("failed to read input from %s: %w", inputName, err)
	}

	// Encrypt the data
	ciphertext, err := aes.Encrypt(key, string(data))
	if err != nil {
		ctx.Logger.Error("Failed to encrypt data", "error", err)
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if s != nil {
		s.Stop()
	}

	// Write output
	if cmd.Output == "" {
		fmt.Println(ciphertext)
	} else {
		// Sanitize output path
		cleanOutPath := filepath.Clean(cmd.Output)
		err := os.WriteFile(cleanOutPath, []byte(ciphertext), 0644)
		if err != nil {
			ctx.Logger.Error("Failed to write output file", "path", cleanOutPath, "error", err)
			return fmt.Errorf("failed to write output file %s: %w", cleanOutPath, err)
		}
		ctx.Logger.Info("File encrypted successfully", "input", inputName, "output", cleanOutPath)
	}

	return nil
}

func (cmd *EncryptCmd) getKey() (string, error) {
	if cmd.Key != "" {
		return cmd.Key, nil
	}

	// Try environment variable
	if envKey := os.Getenv("AES_KEY"); envKey != "" {
		return envKey, nil
	}

	return "", fmt.Errorf("no AES key provided: use --key flag or set AES_KEY environment variable")
}

// DecryptCmd decrypts a file using AES-GCM
type DecryptCmd struct {
	File   string `arg:"" help:"File to decrypt (use '-' for stdin)"`
	Key    string `short:"k" help:"Base64-encoded AES key (if not provided, reads from AES_KEY env var)"`
	Output string `short:"o" help:"Output file (if not provided, prints to stdout)"`
}

func (cmd *DecryptCmd) Run(ctx *CLIContext) error {
	// Get the key
	key, err := cmd.getKey()
	if err != nil {
		ctx.Logger.Error("Failed to get decryption key", "error", err)
		return err
	}

	ctx.Logger.Debug("Decrypting file", "file", cmd.File, "output", cmd.Output)

	// Read input
	var input io.Reader
	var inputName string

	if cmd.File == "-" {
		input = os.Stdin
		inputName = "stdin"
		ctx.Logger.Debug("Reading from stdin")
	} else {
		// Sanitize and validate file path
		cleanPath := filepath.Clean(cmd.File)
		if !strings.HasPrefix(cleanPath, "/") && !strings.HasPrefix(cleanPath, "./") && !strings.HasPrefix(cleanPath, "../") {
			cleanPath = "./" + cleanPath
		}

		file, err := os.Open(cleanPath)
		if err != nil {
			ctx.Logger.Error("Failed to open input file", "path", cleanPath, "error", err)
			return fmt.Errorf("failed to open input file %s: %w", cleanPath, err)
		}
		defer file.Close()

		input = file
		inputName = cleanPath
	}

	// Show spinner for file operations (only if not reading from stdin and not outputting to stdout)
	var s *spinner.Spinner
	if cmd.File != "-" && cmd.Output != "" {
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Decrypting file..."
		s.Start()
		defer s.Stop()
	}

	// Read ciphertext - handle both single line and multiline input
	var ciphertext string
	if cmd.File == "-" {
		// For stdin, read all data and trim whitespace
		data, err := io.ReadAll(input)
		if err != nil {
			ctx.Logger.Error("Failed to read from stdin", "error", err)
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		ciphertext = strings.TrimSpace(string(data))
	} else {
		// For files, try to detect if it's a single line or binary
		scanner := bufio.NewScanner(input)
		if scanner.Scan() {
			ciphertext = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			ctx.Logger.Error("Failed to read input file", "source", inputName, "error", err)
			return fmt.Errorf("failed to read input from %s: %w", inputName, err)
		}
	}

	if ciphertext == "" {
		err := fmt.Errorf("no ciphertext data found in input")
		ctx.Logger.Error("Empty input", "source", inputName, "error", err)
		return err
	}

	// Decrypt the data
	plaintext, err := aes.Decrypt(key, ciphertext)
	if err != nil {
		ctx.Logger.Error("Failed to decrypt data", "error", err)
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	if s != nil {
		s.Stop()
	}

	// Write output
	if cmd.Output == "" {
		fmt.Print(plaintext)
	} else {
		// Sanitize output path
		cleanOutPath := filepath.Clean(cmd.Output)
		err := os.WriteFile(cleanOutPath, []byte(plaintext), 0644)
		if err != nil {
			ctx.Logger.Error("Failed to write output file", "path", cleanOutPath, "error", err)
			return fmt.Errorf("failed to write output file %s: %w", cleanOutPath, err)
		}
		ctx.Logger.Info("File decrypted successfully", "input", inputName, "output", cleanOutPath)
	}

	return nil
}

func (cmd *DecryptCmd) getKey() (string, error) {
	if cmd.Key != "" {
		return cmd.Key, nil
	}

	// Try environment variable
	if envKey := os.Getenv("AES_KEY"); envKey != "" {
		return envKey, nil
	}

	return "", fmt.Errorf("no AES key provided: use --key flag or set AES_KEY environment variable")
}
