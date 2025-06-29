package base64

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty data",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "simple text",
			input:    []byte("hello"),
			expected: "aGVsbG8=",
		},
		{
			name:     "text with special characters",
			input:    []byte("hello world!@#$%^&*()"),
			expected: "aGVsbG8gd29ybGQhQCMkJV4mKigp",
		},
		{
			name:     "binary data",
			input:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expected: "AAECA/8=",
		},
		{
			name:     "unicode text",
			input:    []byte("Hello 世界"),
			expected: "SGVsbG8g5LiW55WM",
		},
		{
			name:     "long text",
			input:    []byte(strings.Repeat("test", 100)),
			expected: "dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdHRlc3R0ZXN0dGVzdA==",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Encode(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEncodeString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "test",
			expected: "dGVzdA==",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			name:     "string with newlines",
			input:    "line1\nline2",
			expected: "bGluZTEKbGluZTI=",
		},
		{
			name:     "string with unicode",
			input:    "café",
			expected: "Y2Fmw6k=",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EncodeString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []byte
		hasError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
			hasError: false,
		},
		{
			name:     "simple encoded text",
			input:    "aGVsbG8=",
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "encoded text with special characters",
			input:    "aGVsbG8gd29ybGQhQCMkJV4mKigp",
			expected: []byte("hello world!@#$%^&*()"),
			hasError: false,
		},
		{
			name:     "encoded binary data",
			input:    "AAECA/8=",
			expected: []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			hasError: false,
		},
		{
			name:     "encoded unicode text",
			input:    "SGVsbG8g5LiW55WM",
			expected: []byte("Hello 世界"),
			hasError: false,
		},
		{
			name:     "input with whitespace",
			input:    "  aGVsbG8=  ",
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "input with tabs and newlines",
			input:    "\t\naGVsbG8=\n\t",
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "invalid base64 - invalid character",
			input:    "aGVsbG8@",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid base64 - wrong padding",
			input:    "aGVsbG8",
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid base64 - incomplete",
			input:    "aGVs",
			expected: []byte("hel"),
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Decode(tc.input)

			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "invalid base64 input")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDecodeToString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "simple encoded text",
			input:    "dGVzdA==",
			expected: "test",
			hasError: false,
		},
		{
			name:     "encoded text with spaces",
			input:    "aGVsbG8gd29ybGQ=",
			expected: "hello world",
			hasError: false,
		},
		{
			name:     "encoded text with newlines",
			input:    "bGluZTEKbGluZTI=",
			expected: "line1\nline2",
			hasError: false,
		},
		{
			name:     "encoded unicode text",
			input:    "Y2Fmw6k=",
			expected: "café",
			hasError: false,
		},
		{
			name:     "input with whitespace",
			input:    "  dGVzdA==  ",
			expected: "test",
			hasError: false,
		},
		{
			name:     "invalid base64",
			input:    "invalid@base64",
			expected: "",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DecodeToString(tc.input)

			if tc.hasError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// TestRoundTrip tests encoding and then decoding data
func TestRoundTrip(t *testing.T) {
	testCases := [][]byte{
		{},
		[]byte("hello"),
		[]byte("hello world!@#$%^&*()"),
		{0x00, 0x01, 0x02, 0x03, 0xFF},
		[]byte(strings.Repeat("test data ", 100)),
	}

	for i, data := range testCases {
		t.Run(fmt.Sprintf("round_trip_%d", i), func(t *testing.T) {
			encoded := Encode(data)
			decoded, err := Decode(encoded)

			require.NoError(t, err)
			assert.Equal(t, data, decoded)
		})
	}
}

// TestRoundTripString tests encoding and then decoding strings
func TestRoundTripString(t *testing.T) {
	testCases := []string{
		"",
		"hello",
		"hello world",
		"line1\nline2\nline3",
		"café with unicode",
		strings.Repeat("test string ", 50),
		"special chars: !@#$%^&*()",
		"tab\tand\nnewline",
	}

	for i, str := range testCases {
		t.Run(fmt.Sprintf("round_trip_string_%d", i), func(t *testing.T) {
			encoded := EncodeString(str)
			decoded, err := DecodeToString(encoded)

			require.NoError(t, err)
			assert.Equal(t, str, decoded)
		})
	}
}

// BenchmarkEncode benchmarks the Encode function
func BenchmarkEncode(b *testing.B) {
	data := []byte(strings.Repeat("benchmark test data ", 100))

	for b.Loop() {
		_ = Encode(data)
	}
}

// BenchmarkDecode benchmarks the Decode function
func BenchmarkDecode(b *testing.B) {
	data := []byte(strings.Repeat("benchmark test data ", 100))
	encoded := Encode(data)

	for b.Loop() {
		_, _ = Decode(encoded)
	}
}

// BenchmarkEncodeString benchmarks the EncodeString function
func BenchmarkEncodeString(b *testing.B) {
	str := strings.Repeat("benchmark test string ", 100)

	for b.Loop() {
		_ = EncodeString(str)
	}
}

// BenchmarkDecodeToString benchmarks the DecodeToString function
func BenchmarkDecodeToString(b *testing.B) {
	str := strings.Repeat("benchmark test string ", 100)
	encoded := EncodeString(str)

	for b.Loop() {
		_, _ = DecodeToString(encoded)
	}
}

// TestEncodeNil tests encoding nil byte slice
func TestEncodeNil(t *testing.T) {
	result := Encode(nil)
	require.Equal(t, "", result)
}

// TestDecodeEdgeCases tests additional edge cases for decoding
func TestDecodeEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "non-printable characters",
			input:       "aGVs\x00bG8=",
			expectedErr: "invalid base64 input",
		},
		{
			name:        "mixed case with invalid chars",
			input:       "aGVsbG8+world",
			expectedErr: "invalid base64 input",
		},
		{
			name:        "string with only whitespace",
			input:       "   \t\n   ",
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Decode(tc.input)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, []byte{}, result)
			}
		})
	}
}

// TestDecodeToStringEdgeCases tests additional edge cases for DecodeToString
func TestDecodeToStringEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "very long invalid input",
			input:       strings.Repeat("invalid@", 1000),
			expectedErr: "invalid base64 input",
		},
		{
			name:        "mixed valid and invalid",
			input:       "aGVsbG8@invalid",
			expectedErr: "invalid base64 input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DecodeToString(tc.input)

			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
			require.Empty(t, result)
		})
	}
}

// TestLargeData tests encoding/decoding of large data
func TestLargeData(t *testing.T) {
	// Test with 1MB of data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	encoded := Encode(largeData)
	require.NotEmpty(t, encoded)

	decoded, err := Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, largeData, decoded)
}

// TestSpecialCharacters tests encoding/decoding of all possible byte values
func TestSpecialCharacters(t *testing.T) {
	// Test all byte values 0-255
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}

	encoded := Encode(allBytes)
	require.NotEmpty(t, encoded)

	decoded, err := Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, allBytes, decoded)
}

// TestConsistencyWithStandardLibrary ensures our functions behave the same as stdlib
func TestConsistencyWithStandardLibrary(t *testing.T) {
	testData := [][]byte{
		{},
		{0},
		{255},
		[]byte("test"),
		[]byte("Hello, 世界!"),
		{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
	}

	for i, data := range testData {
		t.Run(fmt.Sprintf("consistency_%d", i), func(t *testing.T) {
			// Compare our Encode with stdlib
			ourResult := Encode(data)
			stdResult := base64.StdEncoding.EncodeToString(data)
			require.Equal(t, stdResult, ourResult)

			// Compare our Decode with stdlib (when valid)
			if len(ourResult) > 0 || len(data) == 0 {
				decoded, err := Decode(ourResult)
				require.NoError(t, err)
				require.Equal(t, data, decoded)
			}
		})
	}
}
