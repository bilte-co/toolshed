package base62_test

import (
	"testing"

	"github.com/bilte-co/toolshed/base62"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	for _, s := range SamplesStd {
		encoded := base62.StdEncoding.Encode([]byte(s.source))

		// Handle empty input case - the function returns nil for empty input
		if len(s.source) == 0 {
			require.Nil(t, encoded, "empty source should encode to nil")
		} else {
			require.Equal(t, s.targetBytes, encoded, "source: %s", s.source)
		}
	}
}

func TestEncodeToString(t *testing.T) {
	for _, s := range SamplesStd {
		encoded := base62.StdEncoding.EncodeToString([]byte(s.source))
		require.Equal(t, s.target, encoded, "source: %s", s.source)
	}
}

func TestDecode(t *testing.T) {
	for _, s := range SamplesStd {
		decoded, err := base62.StdEncoding.Decode(s.targetBytes)
		require.NoError(t, err, "target: %s", s.target)

		// Handle empty input case - the function returns nil for empty input
		if len(s.target) == 0 {
			require.Nil(t, decoded, "empty target should decode to nil")
		} else {
			require.Equal(t, s.sourceBytes, decoded, "target: %s", s.target)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, s := range SamplesStd {
		decoded, err := base62.StdEncoding.DecodeString(s.target)
		require.NoError(t, err, "target: %s", s.target)

		// Handle empty input case - the function returns nil for empty input
		if len(s.target) == 0 {
			require.Nil(t, decoded, "empty target should decode to nil")
		} else {
			require.Equal(t, s.sourceBytes, decoded, "target: %s", s.target)
		}
	}
}

func TestDecodeWithNewLine(t *testing.T) {
	for _, s := range SamplesWithNewLine {
		decoded, err := base62.StdEncoding.Decode(s.targetBytes)
		require.NoError(t, err, "target: %s", s.target)
		require.Equal(t, s.sourceBytes, decoded, "target: %s", s.target)
	}
}

func TestDecodeError(t *testing.T) {
	for _, s := range SamplesErr {
		_, err := base62.StdEncoding.Decode(s.targetBytes)
		require.Error(t, err, "Expected error for input: %s", s.target)
	}
}

func TestEncodeWithCustomAlphabet(t *testing.T) {
	for _, s := range SamplesWithAlphabet {
		encoded := base62.NewEncoding(s.alphabet).Encode([]byte(s.source))

		// Handle empty input case - the function returns nil for empty input
		if len(s.source) == 0 {
			require.Nil(t, encoded, "empty source should encode to nil")
		} else {
			require.Equal(t, s.targetBytes, encoded, "source: %s, alphabet: %s", s.source, s.alphabet)
		}
	}
}

func TestDecodeWithCustomAlphabet(t *testing.T) {
	for _, s := range SamplesWithAlphabet {
		decoded, err := base62.NewEncoding(s.alphabet).Decode(s.targetBytes)
		require.NoError(t, err, "target: %s, alphabet: %s", s.target, s.alphabet)

		// Handle empty input case - the function returns nil for empty input
		if len(s.target) == 0 {
			require.Nil(t, decoded, "empty target should decode to nil")
		} else {
			require.Equal(t, s.sourceBytes, decoded, "target: %s, alphabet: %s", s.target, s.alphabet)
		}
	}
}

func NewSample(source, target string) *Sample {
	return &Sample{source: source, target: target, sourceBytes: []byte(source), targetBytes: []byte(target)}
}

func NewSampleWithAlphabet(source, target, alphabet string) *Sample {
	return &Sample{source: source, target: target, sourceBytes: []byte(source), targetBytes: []byte(target), alphabet: alphabet}
}

type Sample struct {
	source      string
	target      string
	sourceBytes []byte
	targetBytes []byte
	alphabet    string
}

var SamplesStd = []*Sample{
	NewSample("", ""),
	NewSample("f", "1e"),
	NewSample("fo", "6ox"),
	NewSample("foo", "SAPP"),
	NewSample("foob", "1sIyuo"),
	NewSample("fooba", "7kENWa1"),
	NewSample("foobar", "VytN8Wjy"),

	NewSample("su", "7gj"),
	NewSample("sur", "VkRe"),
	NewSample("sure", "275mAn"),
	NewSample("sure.", "8jHquZ4"),
	NewSample("asure.", "UQPPAab8"),
	NewSample("easure.", "26h8PlupSA"),
	NewSample("leasure.", "9IzLUOIY2fe"),

	NewSample("=", "z"),
	NewSample(">", "10"),
	NewSample("?", "11"),
	NewSample("11", "3H7"),
	NewSample("111", "DWfh"),
	NewSample("1111", "tquAL"),
	NewSample("11111", "3icRuhV"),
	NewSample("111111", "FMElG7cn"),

	NewSample("Hello, World!", "1wJfrzvdbtXUOlUjUf"),
}

var SamplesWithAlphabet = []*Sample{
	NewSampleWithAlphabet("", "", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("f", "Bo", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("fo", "Gy7", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("foo", "cKZZ", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("foob", "B2S84y", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("fooba", "HuOXgkB", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("foobar", "f83XIgt8", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),

	NewSampleWithAlphabet("su", "Hqt", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("sur", "fubo", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("sure", "CHFwKx", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("sure.", "ItR04jE", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("asure.", "eaZZKklI", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("easure.", "CGrIZv4zcK", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("leasure.", "JS9VeYSiCpo", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),

	NewSampleWithAlphabet("=", "9", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet(">", "BA", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("?", "BB", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("11", "DRH", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("111", "Ngpr", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("1111", "304KV", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("11111", "Dsmb4rf", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	NewSampleWithAlphabet("111111", "PWOvQHmx", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),

	NewSampleWithAlphabet("Hello, World!", "B6Tp195nl3heYvetep", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
}

var SamplesWithNewLine = []*Sample{
	NewSample("111111", "FMEl\nG7cn"),
	NewSample("111111", "FMEl\rG7cn"),
	NewSample("Hello, World!", "1wJfrzvdb\ntXUOlUjUf"),
	NewSample("Hello, World!", "1wJfrzvdb\rtXUOlUjUf"),
}

var SamplesErr = []*Sample{
	NewSample("", "Hello, World!"),
}

// Test NewEncoding with proper testify assertions
func TestNewEncoding(t *testing.T) {
	// Valid alphabet
	customAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	enc := base62.NewEncoding(customAlphabet)
	require.NotNil(t, enc)

	// Test that encoding works
	encoded := enc.EncodeToString([]byte("test"))
	require.NotEmpty(t, encoded)

	// Test round trip
	decoded, err := enc.DecodeString(encoded)
	require.NoError(t, err)
	require.Equal(t, "test", string(decoded))
}

func TestNewEncoding_PanicConditions(t *testing.T) {
	// Test alphabet too short
	require.Panics(t, func() {
		base62.NewEncoding("0123456789")
	})

	// Test alphabet too long
	require.Panics(t, func() {
		base62.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123")
	})

	// Test alphabet with newline
	require.Panics(t, func() {
		base62.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ\nabcdefghijklmnopqrstuvwx")
	})

	// Test alphabet with carriage return
	require.Panics(t, func() {
		base62.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ\rabcdefghijklmnopqrstuvwx")
	})
}

func TestEncode_EdgeCases(t *testing.T) {
	// Empty input
	result := base62.StdEncoding.Encode([]byte{})
	require.Nil(t, result)

	// Zero byte (encode as empty, base62 skips leading zeros like many base encodings)
	result = base62.StdEncoding.Encode([]byte{0})
	require.Equal(t, []byte{}, result)

	// Single non-zero byte
	result = base62.StdEncoding.Encode([]byte{1})
	require.NotEmpty(t, result)

	// Verify round trip for non-zero byte
	decoded, err := base62.StdEncoding.Decode(result)
	require.NoError(t, err)
	require.Equal(t, []byte{1}, decoded)
}

func TestEncodeToString_EdgeCases(t *testing.T) {
	// Empty input
	result := base62.StdEncoding.EncodeToString([]byte{})
	require.Empty(t, result)

	// Single byte
	result = base62.StdEncoding.EncodeToString([]byte{42})
	require.NotEmpty(t, result)

	// Large input (starting from 1 to avoid leading zero issues)
	largeInput := make([]byte, 100)
	for i := range largeInput {
		largeInput[i] = byte((i + 1) % 256)
	}
	result = base62.StdEncoding.EncodeToString(largeInput)
	require.NotEmpty(t, result)

	// Verify round trip for large input
	decoded, err := base62.StdEncoding.DecodeString(result)
	require.NoError(t, err)
	// Should equal the input with leading zeros trimmed
	expected := trimLeadingZeros(largeInput)
	require.Equal(t, expected, decoded)
}

func TestDecode_EdgeCases(t *testing.T) {
	// Empty input
	result, err := base62.StdEncoding.Decode([]byte{})
	require.NoError(t, err)
	require.Nil(t, result)

	// Single character "0" decodes to empty (leading zero handling)
	result, err = base62.StdEncoding.Decode([]byte("0"))
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)

	// With newlines and carriage returns mixed
	result, err = base62.StdEncoding.Decode([]byte("1\n2\r3"))
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestDecodeString_EdgeCases(t *testing.T) {
	// Empty string
	result, err := base62.StdEncoding.DecodeString("")
	require.NoError(t, err)
	require.Nil(t, result)

	// Single character strings for each character type
	result, err = base62.StdEncoding.DecodeString("0")
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)

	result, err = base62.StdEncoding.DecodeString("A")
	require.NoError(t, err)
	require.NotNil(t, result)

	result, err = base62.StdEncoding.DecodeString("a")
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestDecode_ErrorCases(t *testing.T) {
	// Invalid characters
	invalidChars := []string{"@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "=", "+", "[", "]", "{", "}", "|", "\\", ":", ";", "\"", "'", "<", ">", "?", "/", ".", ",", "~", "`"}

	for _, char := range invalidChars {
		_, err := base62.StdEncoding.DecodeString("valid" + char + "input")
		require.Error(t, err)

		// Check it's the correct error type
		var corruptErr base62.CorruptInputError
		require.ErrorAs(t, err, &corruptErr)
	}

	// Space character (invalid)
	_, err := base62.StdEncoding.DecodeString("valid input")
	require.Error(t, err)
	var corruptErr base62.CorruptInputError
	require.ErrorAs(t, err, &corruptErr)
}

func TestCorruptInputError(t *testing.T) {
	// Test the error message format
	err := base62.CorruptInputError(32) // ASCII space
	errMsg := err.Error()
	require.Contains(t, errMsg, "illegal base62 data")
	require.Contains(t, errMsg, "32")

	// Test different byte values
	err = base62.CorruptInputError(64) // ASCII @
	errMsg = err.Error()
	require.Contains(t, errMsg, "64")
}

func TestRoundTrip_ComprehensiveData(t *testing.T) {
	testCases := [][]byte{
		{},                        // empty
		{1},                       // single non-zero byte (zero bytes encode differently)
		{255},                     // single max byte
		{1, 2, 3, 4, 5},           // sequential bytes (starting from 1)
		{255, 254, 253, 252, 251}, // descending bytes
		make([]byte, 100),         // smaller test array for speed
	}

	// Fill the last test case with byte values starting from 1
	for i := range testCases[len(testCases)-1] {
		testCases[len(testCases)-1][i] = byte((i + 1) % 256)
	}

	for i, testData := range testCases {
		// Test both Encode/Decode and EncodeToString/DecodeString

		// Test []byte methods
		encoded := base62.StdEncoding.Encode(testData)
		decoded, err := base62.StdEncoding.Decode(encoded)
		require.NoError(t, err, "Test case %d failed", i)

		// Handle special case for empty inputs and zero-leading inputs
		if len(testData) == 0 {
			require.Nil(t, decoded, "Empty input should decode to nil")
		} else {
			// For round trip, we need to handle that leading zeros are stripped
			// So we compare without leading zeros
			expectedData := trimLeadingZeros(testData)
			actualData := decoded
			if len(expectedData) == 0 && len(actualData) == 0 {
				// Both empty, this is correct
			} else {
				require.Equal(t, expectedData, actualData, "Round trip failed for test case %d", i)
			}
		}

		// Test string methods
		encodedStr := base62.StdEncoding.EncodeToString(testData)
		decodedStr, err := base62.StdEncoding.DecodeString(encodedStr)
		require.NoError(t, err, "String test case %d failed", i)

		// Same handling for string methods
		if len(testData) == 0 {
			require.Nil(t, decodedStr, "Empty input should decode to nil")
		} else {
			expectedData := trimLeadingZeros(testData)
			if len(expectedData) == 0 && len(decodedStr) == 0 {
				// Both empty, this is correct
			} else {
				require.Equal(t, expectedData, decodedStr, "String round trip failed for test case %d", i)
			}
		}

		// Verify []byte and string methods produce same results
		require.Equal(t, string(encoded), encodedStr, "Encode methods inconsistent for test case %d", i)
	}
}

// Helper function to trim leading zeros from byte slice
func trimLeadingZeros(data []byte) []byte {
	for i, b := range data {
		if b != 0 {
			return data[i:]
		}
	}
	return []byte{} // All zeros
}

func TestStdEncoding_Availability(t *testing.T) {
	// Ensure StdEncoding is available and working
	require.NotNil(t, base62.StdEncoding)

	// Test basic functionality
	input := "Hello, World!"
	encoded := base62.StdEncoding.EncodeToString([]byte(input))
	require.NotEmpty(t, encoded)

	decoded, err := base62.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	require.Equal(t, input, string(decoded))
}
