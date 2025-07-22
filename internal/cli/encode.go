package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bilte-co/toolshed/base62"
	"github.com/bilte-co/toolshed/base64"
)

// EncodeCmd represents the encode command group
type EncodeCmd struct {
	Encode EncodeTextCmd `cmd:"" help:"Encode text using various encoding schemes"`
	Decode DecodeTextCmd `cmd:"" help:"Decode text using various encoding schemes"`
}

// EncodeTextCmd encodes text using specified encoding
type EncodeTextCmd struct {
	Text     string `arg:"" help:"Text to encode (use '-' to read from stdin)"`
	Encoding string `short:"e" default:"base64" help:"Encoding scheme (base64, base62)"`
}

func (cmd *EncodeTextCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Encoding text", "encoding", cmd.Encoding)

	var input string
	var err error

	// Check if we should read from stdin
	if cmd.Text == "-" {
		input, err = cmd.readStdin()
		if err != nil {
			ctx.Logger.Error("Failed to read from stdin", "error", err)
			return err
		}
	} else {
		input = cmd.Text
	}

	// Set default encoding if empty
	encoding := cmd.Encoding
	if encoding == "" {
		encoding = "base64"
	}

	var result string
	switch strings.ToLower(encoding) {
	case "base64":
		result = base64.EncodeString(input)
	case "base62":
		result = base62.StdEncoding.EncodeToString([]byte(input))
	default:
		err := fmt.Errorf("unsupported encoding: %s (supported: base64, base62)", encoding)
		ctx.Logger.Error("Unsupported encoding", "encoding", encoding)
		return err
	}

	fmt.Println(result)
	ctx.Logger.Info("Text encoded successfully", "encoding", encoding)
	return nil
}

func (cmd *EncodeTextCmd) readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}
	return strings.TrimRight(string(data), "\n\r"), nil
}

// DecodeTextCmd decodes text using specified encoding
type DecodeTextCmd struct {
	Text     string `arg:"" help:"Text to decode (use '-' to read from stdin)"`
	Encoding string `short:"e" default:"base64" help:"Encoding scheme (base64, base62)"`
}

func (cmd *DecodeTextCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Decoding text", "encoding", cmd.Encoding)

	var input string
	var err error

	// Check if we should read from stdin
	if cmd.Text == "-" {
		input, err = cmd.readStdin()
		if err != nil {
			ctx.Logger.Error("Failed to read from stdin", "error", err)
			return err
		}
	} else {
		input = cmd.Text
	}

	// Set default encoding if empty
	encoding := cmd.Encoding
	if encoding == "" {
		encoding = "base64"
	}

	var result string
	switch strings.ToLower(encoding) {
	case "base64":
		result, err = base64.DecodeToString(input)
		if err != nil {
			ctx.Logger.Error("Failed to decode base64", "error", err)
			return fmt.Errorf("failed to decode base64: %w", err)
		}
	case "base62":
		decoded, err := base62.StdEncoding.DecodeString(input)
		if err != nil {
			ctx.Logger.Error("Failed to decode base62", "error", err)
			return fmt.Errorf("failed to decode base62: %w", err)
		}
		result = string(decoded)
	default:
		err := fmt.Errorf("unsupported encoding: %s (supported: base64, base62)", encoding)
		ctx.Logger.Error("Unsupported encoding", "encoding", encoding)
		return err
	}

	fmt.Println(result)
	ctx.Logger.Info("Text decoded successfully", "encoding", encoding)
	return nil
}

func (cmd *DecodeTextCmd) readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}
	return strings.TrimRight(string(data), "\n\r"), nil
}
