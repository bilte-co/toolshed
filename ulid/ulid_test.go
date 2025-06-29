package ulid

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/base62"
	"github.com/stretchr/testify/require"
)

func TestCreateULID_NoPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("", ts)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	ulidBytes, err := base62.StdEncoding.DecodeString(id)
	require.NoError(t, err)
	require.Len(t, ulidBytes, 16)
}

func TestCreateULID_WithPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("test", ts)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(id, "test_"))

	parts := strings.SplitN(id, "_", 2)
	require.Len(t, parts, 2)
	ulidBytes, err := base62.StdEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	require.Len(t, ulidBytes, 16)
}

func TestDecode_ValidULID(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("myobj", ts)
	require.NoError(t, err)

	ulidVal, err := Decode(id)
	require.NoError(t, err)
	require.Equal(t, 16, len(ulidVal))
}

func TestDecode_NoPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("", ts)
	require.NoError(t, err)

	ulidVal, err := Decode(id)
	require.NoError(t, err)
	require.Equal(t, 16, len(ulidVal))
}

func TestDecode_InvalidBase62(t *testing.T) {
	_, err := Decode("prefix_***invalid***")
	require.Error(t, err)
}

func TestDecode_InvalidLength(t *testing.T) {
	_, err := Decode("prefix_" + base62.StdEncoding.EncodeToString([]byte("short")))
	require.ErrorContains(t, err, "invalid ULID length")
}

func TestTimestamp_Extract(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	id, err := CreateULID("ts", now)
	require.NoError(t, err)

	extracted, err := Timestamp(id)
	require.NoError(t, err)

	// Should match to the millisecond
	require.WithinDuration(t, now, extracted, time.Millisecond)
}

func TestTimestamp_InvalidID(t *testing.T) {
	_, err := Timestamp("bad_string")
	require.Error(t, err)
}

func Test_stripPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"abc_def_ghi", "ghi", false},
		{"abc", "abc", false},
		{"no_suffix_", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		out, err := stripPrefix(tt.input)
		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, out)
		}
	}
}

// Additional tests for edge cases and error conditions

func TestCreateULID_EdgeCases(t *testing.T) {
	// Test with very long prefix
	longPrefix := strings.Repeat("a", 100)
	id, err := CreateULID(longPrefix, time.Now())
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(id, longPrefix+"_"))

	// Test with prefix containing underscores
	id, err = CreateULID("prefix_with_underscores", time.Now())
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(id, "prefix_with_underscores_"))

	// Test with special characters in prefix
	id, err = CreateULID("test-123.special", time.Now())
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(id, "test-123.special_"))
}



func TestCreateULID_Uniqueness(t *testing.T) {
	now := time.Now()
	ids := make(map[string]bool)
	
	// Generate multiple ULIDs at the same timestamp
	for i := 0; i < 100; i++ {
		id, err := CreateULID("test", now)
		require.NoError(t, err)
		require.False(t, ids[id], "ULID should be unique")
		ids[id] = true
	}
}

func TestDecode_EmptyString(t *testing.T) {
	_, err := Decode("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ULID")
}

func TestDecode_OnlyPrefix(t *testing.T) {
	_, err := Decode("prefix_")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ULID")
}

func TestDecode_OnlyUnderscore(t *testing.T) {
	_, err := Decode("_")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ULID")
}

func TestDecode_MultipleUnderscores(t *testing.T) {
	// Create a valid ULID first
	ts := time.Now()
	id, err := CreateULID("", ts)
	require.NoError(t, err)
	
	// Create test string with multiple underscores
	testID := "first_second_third_" + id
	
	// Should decode successfully (uses last underscore)
	ulidVal, err := Decode(testID)
	require.NoError(t, err)
	require.Equal(t, 16, len(ulidVal))
}

func TestTimestamp_EmptyString(t *testing.T) {
	_, err := Timestamp("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ULID")
}

func TestTimestamp_OnlyPrefix(t *testing.T) {
	_, err := Timestamp("prefix_")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid ULID")
}

func TestTimestamp_InvalidBase62(t *testing.T) {
	_, err := Timestamp("prefix_***invalid***")
	require.Error(t, err)
}

func TestTimestamp_InvalidLength(t *testing.T) {
	shortData := base62.StdEncoding.EncodeToString([]byte("short"))
	_, err := Timestamp("prefix_" + shortData)
	require.ErrorContains(t, err, "invalid ULID length")
}

func TestTimestamp_BoundaryValues(t *testing.T) {
	// Test with recent past timestamp (within ULID range)
	past := time.Unix(1640995200, 0) // 2022-01-01
	id, err := CreateULID("", past)
	require.NoError(t, err)
	
	extracted, err := Timestamp(id)
	require.NoError(t, err)
	require.WithinDuration(t, past, extracted, time.Millisecond)

	// Test with current time
	now := time.Now()
	id, err = CreateULID("", now)
	require.NoError(t, err)
	
	extracted, err = Timestamp(id)
	require.NoError(t, err)
	require.WithinDuration(t, now, extracted, time.Millisecond)
}

func TestSafeUint64ToInt64(t *testing.T) {
	// Test valid conversion
	result, err := safeUint64ToInt64(12345)
	require.NoError(t, err)
	require.Equal(t, int64(12345), result)

	// Test max valid value
	result, err = safeUint64ToInt64(math.MaxInt64)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), result)

	// Test overflow
	_, err = safeUint64ToInt64(math.MaxInt64 + 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "uint64 value too large for int64")

	// Test max uint64 value
	_, err = safeUint64ToInt64(math.MaxUint64)
	require.Error(t, err)
	require.Contains(t, err.Error(), "uint64 value too large for int64")
}

func TestRoundTrip(t *testing.T) {
	// Test round-trip encoding/decoding
	testCases := []struct {
		prefix string
		time   time.Time
	}{
		{"", time.Now()},
		{"user", time.Now()},
		{"obj_with_undercore", time.Now()},
		{"test", time.Unix(1640995200, 0)}, // 2022-01-01
	}

	for _, tc := range testCases {
		t.Run("prefix="+tc.prefix, func(t *testing.T) {
			// Create ULID
			id, err := CreateULID(tc.prefix, tc.time)
			require.NoError(t, err)

			// Decode it back
			decoded, err := Decode(id)
			require.NoError(t, err)
			require.Equal(t, 16, len(decoded))

			// Extract timestamp
			extractedTime, err := Timestamp(id)
			require.NoError(t, err)
			require.WithinDuration(t, tc.time.Truncate(time.Millisecond), extractedTime, time.Millisecond)
		})
	}
}

func TestStripPrefix_EdgeCases(t *testing.T) {
	// Test single character after underscore
	result, err := stripPrefix("prefix_a")
	require.NoError(t, err)
	require.Equal(t, "a", result)

	// Test multiple consecutive underscores
	result, err = stripPrefix("prefix__suffix")
	require.NoError(t, err)
	require.Equal(t, "suffix", result)

	// Test underscore at beginning
	result, err = stripPrefix("_suffix")
	require.NoError(t, err)
	require.Equal(t, "suffix", result)

	// Test no underscore
	result, err = stripPrefix("noseparator")
	require.NoError(t, err)
	require.Equal(t, "noseparator", result)
}
