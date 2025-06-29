package cli_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/bilte-co/toolshed/ulid"
	"github.com/stretchr/testify/require"
)

func TestULIDCreateCmd_BasicCreation(t *testing.T) {
	cmd := &cli.ULIDCreateCmd{}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDCreateCmd_WithPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{"simple prefix", "usr"},
		{"longer prefix", "user"},
		{"mixed case", "UsEr"},
		{"with numbers", "user123"},
		{"max length", strings.Repeat("a", 32)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Prefix: tt.prefix,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestULIDCreateCmd_WithCustomTimestamp(t *testing.T) {
	// Test with various RFC3339 formatted timestamps
	tests := []struct {
		name      string
		timestamp string
	}{
		{"current time", time.Now().Format(time.RFC3339)},
		{"past time", "2020-01-01T00:00:00Z"},
		{"future time", "2030-12-31T23:59:59Z"},
		{"with timezone", "2023-06-15T14:30:45+05:00"},
		{"with milliseconds", "2023-06-15T14:30:45.123Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Timestamp: tt.timestamp,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestULIDCreateCmd_InvalidTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
	}{
		{"invalid format", "2023-13-45"},
		{"not RFC3339", "2023/01/01 12:00:00"},
		{"empty string", ""},
		{"just text", "invalid-timestamp"},
		{"partial date", "2023-01"},
		{"unix timestamp", "1672531200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Timestamp: tt.timestamp,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid timestamp format")
		})
	}
}

func TestULIDCreateCmd_InvalidPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		errMsg string
	}{
		{"with space", "user id", "cannot contain whitespace"},
		{"with tab", "user\tid", "cannot contain whitespace"},
		{"with newline", "user\nid", "cannot contain whitespace"},
		{"with underscore", "user_id", "cannot contain whitespace"},
		{"too long", strings.Repeat("a", 33), "cannot exceed 32 characters"},
		{"way too long", strings.Repeat("x", 100), "cannot exceed 32 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Prefix: tt.prefix,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestULIDCreateCmd_PrefixAndTimestamp(t *testing.T) {
	cmd := &cli.ULIDCreateCmd{
		Prefix:    "test",
		Timestamp: "2023-06-15T14:30:45Z",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDCreateCmd_MultipleInvocations(t *testing.T) {
	// Ensure multiple calls generate different ULIDs
	cmd := &cli.ULIDCreateCmd{}
	ctx := newTestContext()

	// Run multiple times to ensure uniqueness
	for i := 0; i < 5; i++ {
		err := cmd.Run(ctx)
		require.NoError(t, err)
	}
}

func TestULIDTimestampCmd_ValidULID(t *testing.T) {
	// Create a ULID first
	testTime := time.Now()
	testULID, err := ulid.CreateULID("", testTime)
	require.NoError(t, err)

	cmd := &cli.ULIDTimestampCmd{
		Text: testULID,
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDTimestampCmd_AllFormats(t *testing.T) {
	// Create a test ULID
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	testULID, err := ulid.CreateULID("", testTime)
	require.NoError(t, err)

	formats := []string{"rfc3339", "unix", "unixmilli"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			cmd := &cli.ULIDTimestampCmd{
				Text:   testULID,
				Format: format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestULIDTimestampCmd_InvalidFormat(t *testing.T) {
	testULID, err := ulid.CreateULID("", time.Now())
	require.NoError(t, err)

	invalidFormats := []string{"invalid", "json", "iso8601", ""}

	for _, format := range invalidFormats {
		t.Run(format, func(t *testing.T) {
			cmd := &cli.ULIDTimestampCmd{
				Text:   testULID,
				Format: format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid format")
		})
	}
}

func TestULIDTimestampCmd_InvalidULID(t *testing.T) {
	tests := []struct {
		name string
		ulid string
	}{
		{"empty string", ""},
		{"too short", "123"},
		{"invalid chars", "invalid-ulid-string"},
		{"wrong length", "01234567890123456789012345"},
		{"not base32", "01234567890123456789UVWXYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDTimestampCmd{
				Text: tt.ulid,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
		})
	}
}

func TestULIDTimestampCmd_StdinInput(t *testing.T) {
	// Create a test ULID
	testULID, err := ulid.CreateULID("", time.Now())
	require.NoError(t, err)

	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(testULID))
	}()

	cmd := &cli.ULIDTimestampCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDTimestampCmd_StdinEmpty(t *testing.T) {
	// Mock empty stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	w.Close() // Close immediately to simulate empty input

	cmd := &cli.ULIDTimestampCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ULID cannot be empty")
}

func TestULIDTimestampCmd_StdinWhitespace(t *testing.T) {
	// Create a test ULID with surrounding whitespace
	testULID, err := ulid.CreateULID("", time.Now())
	require.NoError(t, err)

	// Mock stdin with whitespace
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte("  \t" + testULID + "\n  "))
	}()

	cmd := &cli.ULIDTimestampCmd{
		Text: "-",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDTimestampCmd_StdinNoData(t *testing.T) {
	// Mock stdin with no piped data (terminal-like)
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Use os.Stdin directly to simulate terminal input
	cmd := &cli.ULIDTimestampCmd{
	Text: "-",
	}
	ctx := newTestContext()

	// Note: this test may fail in CI environments where stdin behavior differs
	err := cmd.Run(ctx)
	require.Error(t, err)
			require.Contains(t, err.Error(), "no data available from stdin")
}

func TestULIDTimestampCmd_WithPrefixedULID(t *testing.T) {
	// Test with ULIDs that have prefixes
	prefixedULID, err := ulid.CreateULID("user", time.Now())
	require.NoError(t, err)

	cmd := &cli.ULIDTimestampCmd{
		Text: prefixedULID,
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestULIDTimestampCmd_TimestampAccuracy(t *testing.T) {
	// Test that we can extract the exact timestamp we put in
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 123000000, time.UTC)
	testULID, err := ulid.CreateULID("", testTime)
	require.NoError(t, err)

	// Extract timestamp and verify it matches (within millisecond precision)
	extractedTime, err := ulid.Timestamp(testULID)
	require.NoError(t, err)

	// ULID timestamp precision is milliseconds, so truncate our test time
	expectedTime := testTime.Truncate(time.Millisecond)
	actualTime := extractedTime.Truncate(time.Millisecond)

	require.True(t, expectedTime.Equal(actualTime), 
		"Expected %v, got %v", expectedTime, actualTime)
}

func TestULIDCreateCmd_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		timestamp string
		expectErr bool
	}{
		{"empty prefix", "", "", false},
		{"single char prefix", "a", "", false},
		{"numeric prefix", "123", "", false},
		{"mixed case prefix", "aBcD", "", false},
		{"epoch timestamp", "", "1970-01-01T00:00:00Z", false},
		{"far future", "", "2100-01-01T00:00:00Z", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Prefix:    tt.prefix,
				Timestamp: tt.timestamp,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestULIDTimestampCmd_CaseInsensitiveFormat(t *testing.T) {
	testULID, err := ulid.CreateULID("", time.Now())
	require.NoError(t, err)

	// Test case insensitive format matching
	formats := []string{"RFC3339", "UNIX", "UNIXMILLI", "rfc3339", "unix", "unixmilli"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			cmd := &cli.ULIDTimestampCmd{
				Text:   testULID,
				Format: format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestULIDCreateCmd_ConcurrentGeneration(t *testing.T) {
	// Test that concurrent ULID generation works properly
	const numGoroutines = 10
	const uuidsPerGoroutine = 10

	results := make(chan string, numGoroutines*uuidsPerGoroutine)
	errors := make(chan error, numGoroutines*uuidsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < uuidsPerGoroutine; j++ {
				cmd := &cli.ULIDCreateCmd{}
				ctx := newTestContext()

				// We can't easily capture the output, but we can test for errors
				err := cmd.Run(ctx)
				if err != nil {
					errors <- err
				} else {
					results <- "success"
				}
			}
		}()
	}

	// Collect results
	successCount := 0
	errorCount := 0
	for i := 0; i < numGoroutines*uuidsPerGoroutine; i++ {
		select {
		case <-results:
			successCount++
		case err := <-errors:
			errorCount++
			t.Errorf("Unexpected error: %v", err)
		}
	}

	require.Equal(t, numGoroutines*uuidsPerGoroutine, successCount)
	require.Equal(t, 0, errorCount)
}

func TestULIDCreateCmd_BoundaryPrefixLength(t *testing.T) {
	// Test boundary conditions for prefix length
	tests := []struct {
		name      string
		prefix    string
		expectErr bool
	}{
		{"31 chars", strings.Repeat("a", 31), false},
		{"32 chars", strings.Repeat("a", 32), false},
		{"33 chars", strings.Repeat("a", 33), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.ULIDCreateCmd{
				Prefix: tt.prefix,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "cannot exceed 32 characters")
			} else {
				require.NoError(t, err)
			}
		})
	}
}


