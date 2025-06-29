package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheck_WeakPassword(t *testing.T) {
	ok, msg := Check("123456")
	require.False(t, ok)
	require.NotEmpty(t, msg)
	require.Contains(t, msg, "insecure password")
}

func TestCheck_StrongPassword(t *testing.T) {
	ok, msg := Check("G@7e*vS93^8!bdT2")
	require.True(t, ok)
	require.Empty(t, msg)
}

func TestCheckEntropy_HighThreshold(t *testing.T) {
	ok, msg := CheckEntropy("strongPassword123!", 120.0)
	require.False(t, ok)
	require.NotEmpty(t, msg)
}

func TestCheckEntropy_LowThreshold(t *testing.T) {
	ok, msg := CheckEntropy("weakbutpassable", 20.0)
	require.True(t, ok)
	require.Empty(t, msg)
}

func TestCheck_EmptyPassword(t *testing.T) {
	ok, msg := Check("")
	require.False(t, ok)
	require.NotEmpty(t, msg)
}

func TestCheck_RepetitivePassword(t *testing.T) {
	ok, msg := Check("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	require.False(t, ok)
	require.Contains(t, msg, "more special characters")
}

// Additional comprehensive tests

func TestCheck_VariousPasswordStrengths(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"very weak numeric", "123", false},
		{"weak lowercase", "password", false},
		{"medium mixed case", "Password", false},
		{"good mixed with numbers", "Password123", false},
		{"strong with special chars", "Password123!", true},
		{"very strong complex", "MyStr0ng!P@ssw0rd2024", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := Check(tt.password)
			require.Equal(t, tt.expected, ok, "Password: %s, Message: %s", tt.password, msg)
			if !ok {
				require.NotEmpty(t, msg)
			} else {
				require.Empty(t, msg)
			}
		})
	}
}

func TestCheck_BoundaryValues(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"single char", "a", false},
		{"two chars", "ab", false},
		{"very long weak", "a" + string(make([]byte, 1000)), false},
		{"long but strong", "MyStr0nG!P@ssw0rd" + string(make([]byte, 100)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := Check(tt.password)
			require.Equal(t, tt.expected, ok, "Password length: %d, Message: %s", len(tt.password), msg)
		})
	}
}

func TestCheck_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"unicode chars", "Пароль123!", true}, // cyrillic characters provide high entropy
		{"emoji password", "Password123😀", true}, // emoji also provides high entropy
		{"mixed symbols", "P@ssw0rd#$%^&*()", true},
		{"only symbols", "!@#$%^&*()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := Check(tt.password)
			require.Equal(t, tt.expected, ok, "Password: %s, Message: %s", tt.password, msg)
		})
	}
}

func TestCheckEntropy_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		minEntropy  float64
		expected    bool
		description string
	}{
		{"zero entropy", "password", 0.0, true, "should pass with zero entropy requirement"},
		{"negative entropy", "password", -10.0, true, "should pass with negative entropy requirement"},
		{"very high entropy", "MyStr0ng!P@ssw0rd2024", 200.0, false, "should fail with impossibly high entropy requirement"},
		{"exact boundary low", "simplepass", 50.0, false, "should test around typical boundary"},
		{"exact boundary high", "Str0ngP@ssw0rd!", 45.0, true, "should pass reasonable boundary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := CheckEntropy(tt.password, tt.minEntropy)
			require.Equal(t, tt.expected, ok, "%s: Password: %s, Entropy: %f, Message: %s", 
				tt.description, tt.password, tt.minEntropy, msg)
		})
	}
}

func TestCheckEntropy_EmptyPassword(t *testing.T) {
	ok, msg := CheckEntropy("", 30.0)
	require.False(t, ok)
	require.NotEmpty(t, msg)
	require.Contains(t, msg, "insecure password")
}

func TestCheckEntropy_WhitespacePasswords(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"only spaces", "     ", false},
		{"tabs and spaces", "\t\t   \t", false},
		{"newlines", "\n\n\n", false},
		{"mixed whitespace", " \t\n ", false},
		{"password with spaces", "My Str0ng P@ssw0rd!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := CheckEntropy(tt.password, 40.0)
			require.Equal(t, tt.expected, ok, "Password: %q, Message: %s", tt.password, msg)
		})
	}
}

func TestDefaultEntropy_Value(t *testing.T) {
	require.Equal(t, 60.0, DefaultEntropy, "DefaultEntropy constant should be 60.0")
}

func TestCheck_UsesDefaultEntropy(t *testing.T) {
	// Test that Check() uses the same logic as CheckEntropy() with DefaultEntropy
	password := "TestP@ssw0rd123!"
	
	checkResult, checkMsg := Check(password)
	entropyResult, entropyMsg := CheckEntropy(password, DefaultEntropy)
	
	require.Equal(t, checkResult, entropyResult, "Check() should behave same as CheckEntropy() with DefaultEntropy")
	require.Equal(t, checkMsg, entropyMsg, "Error messages should be identical")
}

func TestCheck_CommonWeakPatterns(t *testing.T) {
	weakPatterns := []string{
		"password",
		"Password",
		"password123",
		"Password123",
		"123456789",
		"qwerty",
		"QWERTY",
		"abc123",
		"admin",
		"letmein",
		"welcome",
		"monkey",
		"dragon",
	}

	for _, pattern := range weakPatterns {
		t.Run("weak_"+pattern, func(t *testing.T) {
			ok, msg := Check(pattern)
			require.False(t, ok, "Common weak pattern should fail: %s", pattern)
			require.NotEmpty(t, msg, "Should provide error message for weak pattern: %s", pattern)
		})
	}
}

func TestCheckEntropy_FloatPrecision(t *testing.T) {
	password := "Str0ngP@ssw0rd!"
	
	// Test with various float precision values
	tests := []struct {
		entropy float64
		name    string
	}{
		{45.0, "whole number"},
		{45.5, "one decimal"},
		{45.123456789, "many decimals"},
		{0.1, "very small"},
		{999.999, "very large"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, msg := CheckEntropy(password, tt.entropy)
			// Just ensure it doesn't panic and returns valid results
			require.IsType(t, true, ok)
			require.IsType(t, "", msg)
		})
	}
}
