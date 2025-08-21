package cli_test

import (
	"os"
	"testing"

	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/bilte-co/toolshed/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestBishopStringCmd_Basic(t *testing.T) {
	cmd := &cli.BishopStringCmd{
		Text:      "test",
		Width:     7,
		Height:    5,
		Algorithm: "md5",
	}

	ctx := createTestContext(t)
	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestBishopStringCmd_CustomOptions(t *testing.T) {
	cmd := &cli.BishopStringCmd{
		Text:      "custom",
		Width:     5,
		Height:    3,
		Symbols:   " .o+",
		StartChar: "A",
		EndChar:   "Z",
		NoBorder:  true,
		Algorithm: "sha256",
	}

	ctx := createTestContext(t)
	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestBishopStringCmd_InvalidAlgorithm(t *testing.T) {
	cmd := &cli.BishopStringCmd{
		Text:      "test",
		Algorithm: "invalid",
	}

	ctx := createTestContext(t)
	err := cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid algorithm")
}

func TestBishopFileCmd_Temp(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "bishop_test")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString("test content")
	require.NoError(t, err)
	tmpfile.Close()

	cmd := &cli.BishopFileCmd{
		Path:      tmpfile.Name(),
		Width:     5,
		Height:    3,
		NoBorder:  true,
		Algorithm: "md5",
	}

	ctx := createTestContext(t)
	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestBishopFileCmd_Raw(t *testing.T) {
	// Create a temporary file with specific bytes
	tmpfile, err := os.CreateTemp("", "bishop_test")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte{0x01, 0x02, 0x03, 0x04})
	require.NoError(t, err)
	tmpfile.Close()

	cmd := &cli.BishopFileCmd{
		Path:     tmpfile.Name(),
		Width:    5,
		Height:   3,
		NoBorder: true,
		Raw:      true,
	}

	ctx := createTestContext(t)
	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestBishopFileCmd_NonExistentFile(t *testing.T) {
	cmd := &cli.BishopFileCmd{
		Path: "/nonexistent/file",
	}

	ctx := createTestContext(t)
	err := cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read file")
}

func TestBishopStdinCmd_Raw(t *testing.T) {
	// Simulate stdin
	input := "test input from stdin"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	// Write data to the pipe
	go func() {
		defer w.Close()
		w.WriteString(input)
	}()

	cmd := &cli.BishopStdinCmd{
		Width:     5,
		Height:    3,
		NoBorder:  true,
		Algorithm: "md5",
	}

	ctx := createTestContext(t)
	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestValidateDimensions(t *testing.T) {
	// Valid dimensions
	require.NoError(t, cli.ValidateDimensions(3, 3))
	require.NoError(t, cli.ValidateDimensions(17, 9))
	require.NoError(t, cli.ValidateDimensions(50, 50))

	// Invalid dimensions
	require.Error(t, cli.ValidateDimensions(2, 5))   // width too small
	require.Error(t, cli.ValidateDimensions(5, 2))   // height too small
	require.Error(t, cli.ValidateDimensions(201, 5)) // width too large
	require.Error(t, cli.ValidateDimensions(5, 201)) // height too large
}

func TestParseCharacter(t *testing.T) {
	// Basic characters
	char, err := cli.ParseCharacter("A", 'X')
	require.NoError(t, err)
	require.Equal(t, 'A', char)

	// Empty string should return default
	char, err = cli.ParseCharacter("", 'X')
	require.NoError(t, err)
	require.Equal(t, 'X', char)

	// Escape sequences
	char, err = cli.ParseCharacter("\\n", 'X')
	require.NoError(t, err)
	require.Equal(t, '\n', char)

	char, err = cli.ParseCharacter("\\t", 'X')
	require.NoError(t, err)
	require.Equal(t, '\t', char)

	char, err = cli.ParseCharacter("\\s", 'X')
	require.NoError(t, err)
	require.Equal(t, ' ', char)

	// Numeric codes
	char, err = cli.ParseCharacter("65", 'X')
	require.NoError(t, err)
	require.Equal(t, 'A', char)

	// Invalid cases
	_, err = cli.ParseCharacter("abc", 'X')
	require.Error(t, err)

	_, err = cli.ParseCharacter("256", 'X') // out of range
	require.Error(t, err)
}

func createTestContext(t *testing.T) *cli.CLIContext {
	return testutil.NewTestContext()
}
