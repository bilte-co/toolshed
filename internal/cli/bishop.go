package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bilte-co/toolshed/bishop"
)

// BishopCmd represents the bishop command group
type BishopCmd struct {
	String BishopStringCmd `cmd:"" help:"Generate ASCII art from a string"`
	File   BishopFileCmd   `cmd:"" help:"Generate ASCII art from a file"`
	Stdin  BishopStdinCmd  `cmd:"" help:"Generate ASCII art from stdin"`
}

// BishopStringCmd generates ASCII art from a string
type BishopStringCmd struct {
	Text      string `arg:"" help:"Text to generate ASCII art from"`
	Width     int    `short:"w" default:"17" help:"Grid width (minimum 3)"`
	Height    int    `short:"h" default:"9" help:"Grid height (minimum 3)"`
	Symbols   string `short:"s" help:"Custom symbols for visit counts (e.g., ' .o+=')" `
	StartChar string `long:"start" default:"S" help:"Start position marker"`
	EndChar   string `long:"end" default:"E" help:"End position marker"`
	NoBorder  bool   `short:"b" help:"Hide decorative border"`
	Algorithm string `short:"a" default:"md5" help:"Hash algorithm (md5, sha256)"`
}

func (cmd *BishopStringCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating bishop art from string",
		"width", cmd.Width,
		"height", cmd.Height,
		"algorithm", cmd.Algorithm,
		"noborder", cmd.NoBorder)

	opts, err := cmd.buildOptions()
	if err != nil {
		ctx.Logger.Error("Invalid options", "error", err)
		return err
	}

	var result string
	switch strings.ToLower(cmd.Algorithm) {
	case "md5":
		result = bishop.GenerateFromString(cmd.Text, opts)
	case "sha256":
		result = bishop.GenerateFromStringSHA256(cmd.Text, opts)
	default:
		ctx.Logger.Error("Invalid algorithm", "algorithm", cmd.Algorithm)
		return fmt.Errorf("invalid algorithm '%s' (supported: md5, sha256)", cmd.Algorithm)
	}

	fmt.Print(result)
	ctx.Logger.Info("Bishop art generated successfully")
	return nil
}

func (cmd *BishopStringCmd) buildOptions() (*bishop.Options, error) {
	opts := bishop.DefaultOptions()

	if cmd.Width > 0 {
		opts.Width = cmd.Width
	}
	if cmd.Height > 0 {
		opts.Height = cmd.Height
	}

	if cmd.Symbols != "" {
		opts.Symbols = []rune(cmd.Symbols)
	}

	if len(cmd.StartChar) > 0 {
		opts.StartChar = []rune(cmd.StartChar)[0]
	}

	if len(cmd.EndChar) > 0 {
		opts.EndChar = []rune(cmd.EndChar)[0]
	}

	opts.ShowBorder = !cmd.NoBorder

	return opts, nil
}

// BishopFileCmd generates ASCII art from a file
type BishopFileCmd struct {
	Path      string `arg:"" help:"File path to read from" type:"existingfile"`
	Width     int    `short:"w" default:"17" help:"Grid width (minimum 3)"`
	Height    int    `short:"h" default:"9" help:"Grid height (minimum 3)"`
	Symbols   string `short:"s" help:"Custom symbols for visit counts (e.g., ' .o+=')" `
	StartChar string `long:"start" default:"S" help:"Start position marker"`
	EndChar   string `long:"end" default:"E" help:"End position marker"`
	NoBorder  bool   `short:"b" help:"Hide decorative border"`
	Algorithm string `short:"a" default:"md5" help:"Hash algorithm (md5, sha256)"`
	Raw       bool   `short:"r" help:"Use raw file bytes instead of hashing"`
}

func (cmd *BishopFileCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating bishop art from file",
		"path", cmd.Path,
		"width", cmd.Width,
		"height", cmd.Height,
		"algorithm", cmd.Algorithm,
		"raw", cmd.Raw)

	// Read file
	data, err := os.ReadFile(cmd.Path)
	if err != nil {
		ctx.Logger.Error("Failed to read file", "path", cmd.Path, "error", err)
		return fmt.Errorf("failed to read file: %w", err)
	}

	opts, err := cmd.buildOptions()
	if err != nil {
		ctx.Logger.Error("Invalid options", "error", err)
		return err
	}

	var result string
	if cmd.Raw {
		// Use raw file bytes directly
		result = bishop.GenerateFromBytes(data, opts)
	} else {
		// Hash the file content first
		switch strings.ToLower(cmd.Algorithm) {
		case "md5":
			result = bishop.GenerateFromString(string(data), opts)
		case "sha256":
			result = bishop.GenerateFromStringSHA256(string(data), opts)
		default:
			ctx.Logger.Error("Invalid algorithm", "algorithm", cmd.Algorithm)
			return fmt.Errorf("invalid algorithm '%s' (supported: md5, sha256)", cmd.Algorithm)
		}
	}

	fmt.Print(result)
	ctx.Logger.Info("Bishop art generated successfully from file", "file", cmd.Path)
	return nil
}

func (cmd *BishopFileCmd) buildOptions() (*bishop.Options, error) {
	opts := bishop.DefaultOptions()

	if cmd.Width > 0 {
		opts.Width = cmd.Width
	}
	if cmd.Height > 0 {
		opts.Height = cmd.Height
	}

	if cmd.Symbols != "" {
		opts.Symbols = []rune(cmd.Symbols)
	}

	if len(cmd.StartChar) > 0 {
		opts.StartChar = []rune(cmd.StartChar)[0]
	}

	if len(cmd.EndChar) > 0 {
		opts.EndChar = []rune(cmd.EndChar)[0]
	}

	opts.ShowBorder = !cmd.NoBorder

	return opts, nil
}

// BishopStdinCmd generates ASCII art from stdin
type BishopStdinCmd struct {
	Width     int    `short:"w" default:"17" help:"Grid width (minimum 3)"`
	Height    int    `short:"h" default:"9" help:"Grid height (minimum 3)"`
	Symbols   string `short:"s" help:"Custom symbols for visit counts (e.g., ' .o+=')" `
	StartChar string `long:"start" default:"S" help:"Start position marker"`
	EndChar   string `long:"end" default:"E" help:"End position marker"`
	NoBorder  bool   `short:"b" help:"Hide decorative border"`
	Algorithm string `short:"a" default:"md5" help:"Hash algorithm (md5, sha256)"`
	Raw       bool   `short:"r" help:"Use raw bytes instead of hashing"`
}

func (cmd *BishopStdinCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Generating bishop art from stdin",
		"width", cmd.Width,
		"height", cmd.Height,
		"algorithm", cmd.Algorithm,
		"raw", cmd.Raw)

	// Read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		ctx.Logger.Error("Failed to read from stdin", "error", err)
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	if len(data) == 0 {
		ctx.Logger.Error("No data received from stdin")
		return fmt.Errorf("no data received from stdin")
	}

	opts, err := cmd.buildOptions()
	if err != nil {
		ctx.Logger.Error("Invalid options", "error", err)
		return err
	}

	var result string
	if cmd.Raw {
		// Use raw bytes directly
		result = bishop.GenerateFromBytes(data, opts)
	} else {
		// Hash the input first
		switch strings.ToLower(cmd.Algorithm) {
		case "md5":
			result = bishop.GenerateFromString(string(data), opts)
		case "sha256":
			result = bishop.GenerateFromStringSHA256(string(data), opts)
		default:
			ctx.Logger.Error("Invalid algorithm", "algorithm", cmd.Algorithm)
			return fmt.Errorf("invalid algorithm '%s' (supported: md5, sha256)", cmd.Algorithm)
		}
	}

	fmt.Print(result)
	ctx.Logger.Info("Bishop art generated successfully from stdin")
	return nil
}

func (cmd *BishopStdinCmd) buildOptions() (*bishop.Options, error) {
	opts := bishop.DefaultOptions()

	if cmd.Width > 0 {
		opts.Width = cmd.Width
	}
	if cmd.Height > 0 {
		opts.Height = cmd.Height
	}

	if cmd.Symbols != "" {
		opts.Symbols = []rune(cmd.Symbols)
	}

	if len(cmd.StartChar) > 0 {
		opts.StartChar = []rune(cmd.StartChar)[0]
	}

	if len(cmd.EndChar) > 0 {
		opts.EndChar = []rune(cmd.EndChar)[0]
	}

	opts.ShowBorder = !cmd.NoBorder

	return opts, nil
}

// ValidateDimensions validates width and height parameters
func ValidateDimensions(width, height int) error {
	if width < 3 {
		return fmt.Errorf("width must be at least 3, got %d", width)
	}
	if height < 3 {
		return fmt.Errorf("height must be at least 3, got %d", height)
	}
	if width > 200 {
		return fmt.Errorf("width too large (max 200), got %d", width)
	}
	if height > 200 {
		return fmt.Errorf("height too large (max 200), got %d", height)
	}
	return nil
}

// ParseCharacter safely parses a single character from string
func ParseCharacter(s string, defaultChar rune) (rune, error) {
	if s == "" {
		return defaultChar, nil
	}

	// Handle escape sequences
	if s == "\\n" {
		return '\n', nil
	}
	if s == "\\t" {
		return '\t', nil
	}
	if s == "\\r" {
		return '\r', nil
	}
	if s == "\\s" {
		return ' ', nil
	}

	// Handle numeric codes (e.g., "65" for 'A')
	if num, err := strconv.Atoi(s); err == nil && num >= 32 && num <= 126 {
		return rune(num), nil
	}

	runes := []rune(s)
	if len(runes) != 1 {
		return defaultChar, fmt.Errorf("expected single character, got %q", s)
	}

	return runes[0], nil
}
