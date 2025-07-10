package cli_test

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/bilte-co/toolshed/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestHaikuGenerateCmd_DefaultGeneration(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		Delim: "-", // Set explicit delimiter since default is empty
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_WithCustomToken(t *testing.T) {
	tests := []struct {
		name  string
		token int64
	}{
		{"small token", 10},
		{"medium token", 1000},
		{"large token", 999999},
		{"token of 1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Token: tt.token,
				Delim: "-", // Set explicit delimiter
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHaikuGenerateCmd_WithCustomDelimiter(t *testing.T) {
	tests := []struct {
		name  string
		delim string
	}{
		{"period", "."},
		{"underscore", "_"},
		{"comma", ","},
		{"colon", ":"},
		{"pipe", "|"},
		{"space", " "},
		{"double dash", "--"},
		{"mixed", "._"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Delim: tt.delim,
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHaikuGenerateCmd_NoTokenFlag(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		NoToken: true,
		Delim:   "-", // Set explicit delimiter
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_ZeroToken(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		Token: 0,
		Delim: "-", // Set explicit delimiter
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_NoTokenWithCustomDelim(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		NoToken: true,
		Delim:   ".",
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_TokenAndDelimCombination(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		Token: 100,
		Delim: "_",
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_UnsafeDelimiter(t *testing.T) {
	tests := []struct {
		name  string
		delim string
	}{
		{"empty delimiter", ""},
		{"unsafe characters", "!@#"},
		{"too long", "------"},
		{"non-ASCII", "αβγ"},
		{"script tag", "<script>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Delim: tt.delim,
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "failed to generate haiku")
		})
	}
}

func TestHaikuGenerateCmd_NoTokenPrecedence(t *testing.T) {
	// Test that NoToken flag takes precedence over Token value
	cmd := &cli.HaikuGenerateCmd{
		Token:   9999,
		NoToken: true,
		Delim:   "-",
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_NegativeToken(t *testing.T) {
	cmd := &cli.HaikuGenerateCmd{
		Token: -1,
		Delim: "-", // Set explicit delimiter
	}
	ctx := testutil.NewTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHaikuGenerateCmd_CaptureOutput(t *testing.T) {
	// Capture stdout to verify haiku is printed
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdout = w

	cmd := &cli.HaikuGenerateCmd{
		Token: 100,
		Delim: "-",
	}
	ctx := testutil.NewTestContext()

	go func() {
		defer w.Close()
		err := cmd.Run(ctx)
		require.NoError(t, err)
	}()

	output, err := io.ReadAll(r)
	require.NoError(t, err)

	outputStr := strings.TrimSpace(string(output))
	require.NotEmpty(t, outputStr, "Should output a haiku")

	// Verify it contains delimiter
	require.Contains(t, outputStr, "-", "Should contain delimiter")

	// Verify it has parts (at least 3 for adjective-action-noun, maybe 4 with token)
	parts := strings.Split(outputStr, "-")
	require.GreaterOrEqual(t, len(parts), 3, "Should have at least 3 parts")

	// If it has 4 parts, last should be a number
	if len(parts) == 4 {
		lastPart := parts[len(parts)-1]
		_, err := strconv.ParseInt(lastPart, 10, 64)
		require.NoError(t, err, "Last part should be a number when token is present")
	}
}

func TestHaikuGenerateCmd_CaptureOutputNoToken(t *testing.T) {
	// Capture stdout to verify haiku without token
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdout = w

	cmd := &cli.HaikuGenerateCmd{
		NoToken: true,
		Delim:   ".",
	}
	ctx := testutil.NewTestContext()

	go func() {
		defer w.Close()
		err := cmd.Run(ctx)
		require.NoError(t, err)
	}()

	output, err := io.ReadAll(r)
	require.NoError(t, err)

	outputStr := strings.TrimSpace(string(output))
	require.NotEmpty(t, outputStr, "Should output a haiku")

	// Verify it contains delimiter
	require.Contains(t, outputStr, ".", "Should contain delimiter")

	// Verify no numeric token at the end
	parts := strings.Split(outputStr, ".")
	require.GreaterOrEqual(t, len(parts), 3, "Should have at least 3 parts")

	lastPart := parts[len(parts)-1]
	_, err = strconv.ParseInt(lastPart, 10, 64)
	require.Error(t, err, "Last part should not be a number when NoToken is true")
}

func TestHaikuGenerateCmd_BoundaryDelimiterLength(t *testing.T) {
	tests := []struct {
		name      string
		delim     string
		expectErr bool
	}{
		{"single char", "-", false},
		{"2 chars", "--", false},
		{"3 chars", "---", false},
		{"4 chars", "----", false},
		{"5 chars (max)", "-----", false},
		{"6 chars (too long)", "------", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Delim: tt.delim,
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHaikuGenerateCmd_LargeTokenValues(t *testing.T) {
	tests := []struct {
		name  string
		token int64
	}{
		{"very large", 999999999},
		{"max int32", 2147483647},
		{"larger than default", 99999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Token: tt.token,
				Delim: "-", // Set explicit delimiter
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHaikuGenerateCmd_MultipleInvocations(t *testing.T) {
	// Test that multiple invocations work correctly
	cmd := &cli.HaikuGenerateCmd{
		Delim: "-", // Set explicit delimiter
	}
	ctx := testutil.NewTestContext()

	for i := 0; i < 5; i++ {
		err := cmd.Run(ctx)
		require.NoError(t, err)
	}
}

func TestHaikuGenerateCmd_ConcurrentExecution(t *testing.T) {
	// Test concurrent execution
	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			cmd := &cli.HaikuGenerateCmd{
				Delim: "-", // Set explicit delimiter
			}
			ctx := testutil.NewTestContext()
			errors <- cmd.Run(ctx)
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		require.NoError(t, err)
	}
}

func TestHaikuGenerateCmd_EdgeCaseCombinations(t *testing.T) {
	tests := []struct {
		name    string
		token   int64
		delim   string
		noToken bool
	}{
		{"zero token with no token flag", 0, "-", true},
		{"negative token with custom delim", -100, ".", false},
		{"large token with space delim", 999999, " ", false},
		{"no token with colon", 0, ":", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HaikuGenerateCmd{
				Token:   tt.token,
				Delim:   tt.delim,
				NoToken: tt.noToken,
			}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}
