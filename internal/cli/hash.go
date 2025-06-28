package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/bilte-co/toolshed/hash"
)

// HashCmd represents the hash command group
type HashCmd struct {
	String   HashStringCmd `cmd:"" help:"Hash a string"`
	File     HashFileCmd   `cmd:"" help:"Hash a file"`
	Dir      HashDirCmd    `cmd:"" help:"Hash a directory"`
	HMAC     HMACCmd       `cmd:"" help:"Compute HMAC of data"`
	Validate ValidateCmd   `cmd:"" help:"Validate file against expected hash"`
	Compare  CompareCmd    `cmd:"" help:"Compare two hashes using constant-time comparison"`
}

// HashStringCmd hashes a string
type HashStringCmd struct {
	Text   string `arg:"" help:"Text to hash"`
	Algo   string `short:"a" default:"sha256" help:"Hash algorithm (md5, sha1, sha256, sha512, blake2b)"`
	Format string `short:"f" default:"hex" help:"Output format (hex, base64, raw)"`
	Prefix bool   `short:"p" help:"Prefix output with algorithm name"`
}

func (cmd *HashStringCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Hashing string", "algorithm", cmd.Algo, "format", cmd.Format)

	opts := hash.Options{
		Format: hash.Format(cmd.Format),
		Prefix: cmd.Prefix,
	}

	result, err := hash.HashStringWithOptions(cmd.Text, cmd.Algo, opts)
	if err != nil {
		ctx.Logger.Error("Failed to hash string", "error", err)
		return err
	}

	fmt.Println(result)
	ctx.Logger.Info("Hash computed successfully")
	return nil
}

// HashFileCmd hashes a file
type HashFileCmd struct {
	Path   string `arg:"" help:"File path to hash" type:"existingfile"`
	Algo   string `short:"a" default:"sha256" help:"Hash algorithm (md5, sha1, sha256, sha512, blake2b)"`
	Format string `short:"f" default:"hex" help:"Output format (hex, base64, raw)"`
	Prefix bool   `short:"p" help:"Prefix output with algorithm name"`
}

func (cmd *HashFileCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Hashing file", "path", cmd.Path, "algorithm", cmd.Algo)

	// Check if we should read from stdin
	if cmd.Path == "-" {
		return cmd.hashStdin(ctx)
	}

	// Sanitize path
	cleanPath := filepath.Clean(cmd.Path)
	if !strings.HasPrefix(cleanPath, "/") && !strings.HasPrefix(cleanPath, "./") && !strings.HasPrefix(cleanPath, "../") {
		cleanPath = "./" + cleanPath
	}

	// Show spinner for large files
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Computing hash..."
	s.Start()
	defer s.Stop()

	opts := hash.Options{
		Format: hash.Format(cmd.Format),
		Prefix: cmd.Prefix,
	}

	result, err := hash.HashFileWithOptions(cleanPath, cmd.Algo, opts)
	if err != nil {
		ctx.Logger.Error("Failed to hash file", "path", cleanPath, "error", err)
		return err
	}

	s.Stop()
	fmt.Println(result)
	ctx.Logger.Info("Hash computed successfully", "file", cleanPath)
	return nil
}

func (cmd *HashFileCmd) hashStdin(ctx *CLIContext) error {
	ctx.Logger.Debug("Reading from stdin")

	opts := hash.Options{
		Format: hash.Format(cmd.Format),
		Prefix: cmd.Prefix,
	}

	result, err := hash.HashReaderWithOptions(os.Stdin, cmd.Algo, opts)
	if err != nil {
		ctx.Logger.Error("Failed to hash stdin", "error", err)
		return err
	}

	fmt.Println(result)
	ctx.Logger.Info("Hash computed successfully from stdin")
	return nil
}

// HashDirCmd hashes a directory
type HashDirCmd struct {
	Path      string `arg:"" help:"Directory path to hash" type:"existingdir"`
	Algo      string `short:"a" default:"sha256" help:"Hash algorithm (md5, sha1, sha256, sha512, blake2b)"`
	Format    string `short:"f" default:"hex" help:"Output format (hex, base64, raw)"`
	Prefix    bool   `short:"p" help:"Prefix output with algorithm name"`
	Recursive bool   `short:"r" default:"true" help:"Hash directories recursively"`
}

func (cmd *HashDirCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Hashing directory", "path", cmd.Path, "recursive", cmd.Recursive, "algorithm", cmd.Algo)

	// Sanitize path
	cleanPath := filepath.Clean(cmd.Path)

	// Show spinner for directory operations
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Computing directory hash..."
	s.Start()
	defer s.Stop()

	opts := hash.Options{
		Format: hash.Format(cmd.Format),
		Prefix: cmd.Prefix,
	}

	result, err := hash.HashDirWithOptions(cleanPath, cmd.Algo, cmd.Recursive, opts)
	if err != nil {
		ctx.Logger.Error("Failed to hash directory", "path", cleanPath, "error", err)
		return err
	}

	s.Stop()
	fmt.Println(result)
	ctx.Logger.Info("Directory hash computed successfully", "path", cleanPath)
	return nil
}

// HMACCmd computes HMAC
type HMACCmd struct {
	Text   string `arg:"" help:"Text to compute HMAC for"`
	Key    string `short:"k" required:"" help:"HMAC key"`
	Algo   string `short:"a" default:"sha256" help:"Hash algorithm (md5, sha1, sha256, sha512, blake2b)"`
	Format string `short:"f" default:"hex" help:"Output format (hex, base64, raw)"`
	Prefix bool   `short:"p" help:"Prefix output with algorithm name"`
}

func (cmd *HMACCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Computing HMAC", "algorithm", cmd.Algo, "format", cmd.Format)

	opts := hash.Options{
		Format: hash.Format(cmd.Format),
		Prefix: cmd.Prefix,
	}

	result, err := hash.HMACWithOptions([]byte(cmd.Text), []byte(cmd.Key), cmd.Algo, opts)
	if err != nil {
		ctx.Logger.Error("Failed to compute HMAC", "error", err)
		return err
	}

	fmt.Println(result)
	ctx.Logger.Info("HMAC computed successfully")
	return nil
}

// ValidateCmd validates a file against expected hash
type ValidateCmd struct {
	File     string `arg:"" help:"File to validate" type:"existingfile"`
	Expected string `short:"e" required:"" help:"Expected hash value"`
	Algo     string `short:"a" default:"sha256" help:"Hash algorithm (md5, sha1, sha256, sha512, blake2b)"`
}

func (cmd *ValidateCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Validating file", "file", cmd.File, "algorithm", cmd.Algo)

	// Sanitize path
	cleanPath := filepath.Clean(cmd.File)

	// Show spinner during validation
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Validating file..."
	s.Start()
	defer s.Stop()

	err := hash.ValidateFileChecksum(cleanPath, cmd.Expected, cmd.Algo)
	s.Stop()

	if err != nil {
		ctx.Logger.Error("Validation failed", "error", err)
		return err
	}

	ctx.Logger.Info("Validation successful", "file", cleanPath)
	fmt.Println("✓ Validation successful")
	return nil
}

// CompareCmd compares two hashes using constant-time comparison
type CompareCmd struct {
	Hash1 string `arg:"" help:"First hash to compare"`
	Hash2 string `arg:"" help:"Second hash to compare"`
}

func (cmd *CompareCmd) Run(ctx *CLIContext) error {
	ctx.Logger.Debug("Comparing hashes")

	// Convert hex strings to bytes
	bytes1, err := hex.DecodeString(strings.TrimSpace(cmd.Hash1))
	if err != nil {
		ctx.Logger.Error("Invalid first hash format", "error", err)
		return fmt.Errorf("invalid first hash format: %w", err)
	}

	bytes2, err := hex.DecodeString(strings.TrimSpace(cmd.Hash2))
	if err != nil {
		ctx.Logger.Error("Invalid second hash format", "error", err)
		return fmt.Errorf("invalid second hash format: %w", err)
	}

	equal := hash.EqualConstantTime(bytes1, bytes2)

	if equal {
		ctx.Logger.Info("Hashes match")
		fmt.Println("✓ Hashes are equal")
	} else {
		ctx.Logger.Info("Hashes do not match")
		fmt.Println("✗ Hashes are different")
	}

	return nil
}
