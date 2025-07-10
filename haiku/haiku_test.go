package haiku

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestDefaultReturnsTwoWordsAndInt(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.Haikunate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, "-")

	if len(parts) < 4 {
		t.Errorf("Generated haiku [%s] should have at least 4 parts: %v\n", haiku, parts)
	}

	// Last part should be the token
	lastPart := parts[len(parts)-1]
	_, err = strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		t.Error("Last part is not integer: ", lastPart)
	}
}

func TestNonDefaultReturnsTwoWordsAndInt(t *testing.T) {
	h := Haikunator{delim: ".", token: 1000}

	haiku, err := h.Haikunate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, ".")

	if len(parts) < 4 {
		t.Errorf("Generated haiku [%s] should have at least 4 parts: %v\n", haiku, parts)
	}

	// Last part should be the token
	lastPart := parts[len(parts)-1]
	token, err := strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		t.Error("Last part is not integer: ", lastPart)
	}

	if token < 0 || token > 1000 {
		t.Error("Generated token is outside of bounds 0-1000", token)
	}
}

func TestTokenHaikunateBetweenZeroAndMax(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.TokenHaikunate(1000)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, "-")

	if len(parts) < 4 {
		t.Errorf("Generated haiku [%s] should have at least 4 parts: %v\n", haiku, parts)
	}

	// Last part should be the token
	lastPart := parts[len(parts)-1]
	token, err := strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		t.Error("Last part is not integer: ", lastPart)
	}

	if token < 0 || token > 1000 {
		t.Error("Generated token is outside of bounds 0-1000", token)
	}
}

func TestZeroTokenGeneratesNoTokenHaiku(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.TokenHaikunate(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, "-")

	if len(parts) < 3 {
		t.Errorf("Generated haiku has invalid number of parts, expected at least 3: %v\n", parts)
	}

	// Check that last part is NOT a number (since no token should be added)
	lastPart := parts[len(parts)-1]
	if _, err := strconv.ParseInt(lastPart, 10, 64); err == nil {
		t.Errorf("Zero token should not generate numeric token, but got: %s", lastPart)
	}
}

func TestSpaceDelimHaikuHasCorrectDelimAndNoToken(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.DelimHaikunate(" ")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, " ")

	if len(parts) < 3 {
		t.Errorf("Generated haiku has invalid number of parts, expected at least 3: %v\n", parts)
	}

	// Check that last part is NOT a number (since no token should be added)
	lastPart := parts[len(parts)-1]
	if _, err := strconv.ParseInt(lastPart, 10, 64); err == nil {
		t.Errorf("DelimHaikunate should not generate numeric token, but got: %s", lastPart)
	}
}

func TestTokenDelimHaikuHasDelimAndToken(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.TokenDelimHaikunate(1000, ".")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, ".")

	if len(parts) < 4 {
		t.Errorf("Generated haiku has invalid number of parts, expected at least 4: %v\n", parts)
	}

	// Last part should be the token
	lastPart := parts[len(parts)-1]
	token, err := strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		t.Error("Last part is not integer: ", lastPart)
	}

	if token < 0 || token > 1000 {
		t.Error("Generated token is outside of bounds 0-1000", token)
	}
}

func TestZeroTokenDelimHaikuHasDelimAndNoToken(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.TokenDelimHaikunate(0, " ")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	parts := strings.Split(haiku, " ")

	if len(parts) < 3 {
		t.Errorf("Generated haiku has invalid number of parts, expected at least 3: %v\n", parts)
	}

	// Check that last part is NOT a number (since no token should be added)
	lastPart := parts[len(parts)-1]
	if _, err := strconv.ParseInt(lastPart, 10, 64); err == nil {
		t.Errorf("Zero token should not generate numeric token, but got: %s", lastPart)
	}
}

func TestUnsafeDelimiterReturnsError(t *testing.T) {
	h := NewHaikunator()

	unsafeDelimiters := []string{
		"",           // empty string
		"<script>",   // contains unsafe characters
		"!@#$%^&*()", // contains unsafe characters
		"abcdef",     // too long
		"£€¥",        // non-ASCII characters
		"αβγ",        // non-ASCII characters
	}

	for _, delim := range unsafeDelimiters {
		_, err := h.DelimHaikunate(delim)
		if err == nil {
			t.Errorf("Expected error for unsafe delimiter: %s", delim)
		}

		_, err = h.TokenDelimHaikunate(1000, delim)
		if err == nil {
			t.Errorf("Expected error for unsafe delimiter in TokenDelimHaikunate: %s", delim)
		}
	}
}

func TestSafeDelimitersWork(t *testing.T) {
	h := NewHaikunator()

	safeDelimiters := []string{
		"-", ".", "_", ",", ":", "|", " ",
		"--", "..", "__", ",,", "::", "||", "  ",
		"-_", "._", ",-", " - ", ":|:",
	}

	for _, delim := range safeDelimiters {
		haiku, err := h.DelimHaikunate(delim)
		if err != nil {
			t.Errorf("Unexpected error for safe delimiter %s: %v", delim, err)
		}
		if !strings.Contains(haiku, delim) {
			t.Errorf("Generated haiku does not contain delimiter %s: %s", delim, haiku)
		}
	}
}

func TestHaikuStructureConsistency(t *testing.T) {
	h := NewHaikunator()

	for i := 0; i < 10; i++ {
		haiku, err := h.Haikunate()
		if err != nil {
			t.Errorf("Unexpected error generating haiku: %v", err)
		}

		parts := strings.Split(haiku, "-")
		if len(parts) < 4 {
			t.Errorf("Haiku should have at least 4 parts, got %d: %s", len(parts), haiku)
		}

		// Check that first two parts are adjective and action
		adjective := parts[0]
		action := parts[1]
		lastPart := parts[len(parts)-1]

		if !slices.Contains(ADJECTIVES, adjective) {
			t.Errorf("First part should be an adjective: %s", adjective)
		}
		if !slices.Contains(ACTIONS, action) {
			t.Errorf("Second part should be an action: %s", action)
		}
		// Last part should be a token (integer) for default Haikunate
		if _, err := strconv.ParseInt(lastPart, 10, 64); err != nil {
			t.Errorf("Last part should be a token: %s", lastPart)
		}
	}
}

func TestNegativeTokenBehavior(t *testing.T) {
	h := NewHaikunator()

	haiku, err := h.TokenHaikunate(-1)
	if err != nil {
		t.Errorf("Unexpected error for negative token: %v", err)
	}

	parts := strings.Split(haiku, "-")
	if len(parts) < 3 {
		t.Errorf("Negative token should generate haiku without token, got %d parts: %s", len(parts), haiku)
	}

	// Check that last part is NOT a number (since no token should be added)
	lastPart := parts[len(parts)-1]
	if _, err := strconv.ParseInt(lastPart, 10, 64); err == nil {
		t.Errorf("Negative token should not generate numeric token, but got: %s", lastPart)
	}
}

func TestWordListsNotEmpty(t *testing.T) {
	if len(ADJECTIVES) == 0 {
		t.Error("ADJECTIVES list should not be empty")
	}
	if len(ACTIONS) == 0 {
		t.Error("ACTIONS list should not be empty")
	}
	if len(NOUNS) == 0 {
		t.Error("NOUNS list should not be empty")
	}
}

func TestHaikuRandomness(t *testing.T) {
	h1 := NewHaikunator()
	h2 := NewHaikunator()

	// Generate multiple haikus and check they're different
	haikus := make(map[string]bool)
	for range 20 {
		haiku1, _ := h1.Haikunate()
		haiku2, _ := h2.Haikunate()

		haikus[haiku1] = true
		haikus[haiku2] = true
	}

	// With 40 haikus generated, we should have some variety
	if len(haikus) < 10 {
		t.Errorf("Expected more variety in generated haikus, got %d unique haikus", len(haikus))
	}
}

func TestTokenBoundaries(t *testing.T) {
	h := NewHaikunator()

	// Test with token value 1
	haiku, err := h.TokenHaikunate(1)
	if err != nil {
		t.Errorf("Unexpected error for token 1: %v", err)
	}

	parts := strings.Split(haiku, "-")
	if len(parts) < 4 {
		t.Errorf("Token haiku should have at least 4 parts, got %d: %s", len(parts), haiku)
	}

	// Last part should be the token
	lastPart := parts[len(parts)-1]
	token, err := strconv.ParseInt(lastPart, 10, 64)
	if err != nil {
		t.Errorf("Last part should be a valid integer: %s", lastPart)
	}

	if token != 0 {
		t.Errorf("Token with max 1 should be 0, got %d", token)
	}
}

func TestDelimiterLengthLimits(t *testing.T) {
	h := NewHaikunator()

	// Test maximum allowed delimiter length (5 characters)
	longDelim := "-----"
	_, err := h.DelimHaikunate(longDelim)
	if err != nil {
		t.Errorf("Delimiter of length 5 should be allowed: %v", err)
	}

	// Test too long delimiter (6 characters)
	tooLongDelim := "------"
	_, err = h.DelimHaikunate(tooLongDelim)
	if err == nil {
		t.Error("Delimiter of length 6 should not be allowed")
	}
}
