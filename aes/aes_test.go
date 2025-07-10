package aes_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/bilte-co/toolshed/aes"
	"github.com/stretchr/testify/require"
)

func TestGenerateAESKey_ValidLengths(t *testing.T) {
	validBits := []int{128, 192, 256}
	for _, bits := range validBits {
		key, err := aes.GenerateAESKey(bits)
		require.NoError(t, err)
		require.NotEmpty(t, key)

		decoded, err := base64.StdEncoding.DecodeString(key)
		require.NoError(t, err)
		require.Len(t, decoded, bits/8)
	}
}

func TestGenerateAESKey_InvalidLength(t *testing.T) {
	_, err := aes.GenerateAESKey(100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "AES key length")
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	plaintext := "Hello, secure world!"
	ciphertext, err := aes.Encrypt(key, plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, ciphertext)

	decrypted, err := aes.Decrypt(key, ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	key, err := aes.GenerateAESKey(128)
	require.NoError(t, err)

	ciphertext, err := aes.Encrypt(key, "")
	require.NoError(t, err)

	decrypted, err := aes.Decrypt(key, ciphertext)
	require.NoError(t, err)
	require.Equal(t, "", decrypted)
}

func TestEncrypt_InvalidBase64Key(t *testing.T) {
	_, err := aes.Encrypt("not@@base64", "plaintext")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base64 key")
}

func TestDecrypt_InvalidBase64Ciphertext(t *testing.T) {
	key, _ := aes.GenerateAESKey(128)
	_, err := aes.Decrypt(key, "bad$$base64")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base64 ciphertext")
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	key, _ := aes.GenerateAESKey(128)
	invalidCiphertext := base64.StdEncoding.EncodeToString([]byte("short"))

	_, err := aes.Decrypt(key, invalidCiphertext)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ciphertext too short")
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key, _ := aes.GenerateAESKey(128)

	msg := "top secret"
	ciphertext, err := aes.Encrypt(key, msg)
	require.NoError(t, err)

	// Tamper with ciphertext
	tampered := ciphertext[:len(ciphertext)-2] + "aa"

	_, err = aes.Decrypt(key, tampered)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt")
}

func TestEncrypt_InvalidDecodedKeyLength(t *testing.T) {
	// 10 bytes (invalid for AES)
	invalidKey := base64.StdEncoding.EncodeToString([]byte("1234567890"))
	_, err := aes.Encrypt(invalidKey, "plaintext")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create AES cipher")
}

func TestEncryptDecrypt_VariousInputs(t *testing.T) {
	key, _ := aes.GenerateAESKey(256)
	inputs := []string{
		"short",
		strings.Repeat("a", 1000),
		"!@#$%^&*()_+=-{}[]:\";'<>?,./",
		"",
		"1234567890",
	}

	for _, input := range inputs {
		ciphertext, err := aes.Encrypt(key, input)
		require.NoError(t, err)

		plaintext, err := aes.Decrypt(key, ciphertext)
		require.NoError(t, err)
		require.Equal(t, input, plaintext)
	}
}

func TestGenerateAESKey_MultipleInvocations(t *testing.T) {
	// Test that multiple calls generate different keys
	key1, err := aes.GenerateAESKey(256)
	require.NoError(t, err)
	key2, err := aes.GenerateAESKey(256)
	require.NoError(t, err)
	require.NotEqual(t, key1, key2)
}

func TestGenerateAESKey_ZeroAndNegativeValues(t *testing.T) {
	testCases := []int{0, -1, -128, -256}
	for _, bits := range testCases {
		_, err := aes.GenerateAESKey(bits)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid AES key length")
	}
}

func TestGenerateAESKey_BoundaryValues(t *testing.T) {
	// Test values around valid lengths
	invalidValues := []int{127, 129, 191, 193, 255, 257}
	for _, bits := range invalidValues {
		_, err := aes.GenerateAESKey(bits)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid AES key length")
	}
}

func TestEncrypt_EmptyKey(t *testing.T) {
	_, err := aes.Encrypt("", "plaintext")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create AES cipher")
}

func TestDecrypt_EmptyKey(t *testing.T) {
	_, err := aes.Decrypt("", "ciphertext")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base64 ciphertext")
}

func TestDecrypt_EmptyCiphertext(t *testing.T) {
	key, _ := aes.GenerateAESKey(128)
	_, err := aes.Decrypt(key, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "ciphertext too short")
}

func TestEncryptDecrypt_UnicodeContent(t *testing.T) {
	key, _ := aes.GenerateAESKey(256)
	unicodeTexts := []string{
		"🔐 Secret emoji message 🗝️",
		"Привет мир",
		"こんにちは世界",
		"مرحبا بالعالم",
		"Mixed: Hello 世界 🌍",
		"\x00\x01\x02\x03", // binary data
	}

	for _, text := range unicodeTexts {
		ciphertext, err := aes.Encrypt(key, text)
		require.NoError(t, err)

		decrypted, err := aes.Decrypt(key, ciphertext)
		require.NoError(t, err)
		require.Equal(t, text, decrypted)
	}
}

func TestEncryptDecrypt_AllKeyLengths(t *testing.T) {
	keyLengths := []int{128, 192, 256}
	plaintext := "Test message for all key lengths"

	for _, bits := range keyLengths {
		key, err := aes.GenerateAESKey(bits)
		require.NoError(t, err)

		ciphertext, err := aes.Encrypt(key, plaintext)
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)

		decrypted, err := aes.Decrypt(key, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, decrypted)
	}
}

func TestEncryptDecrypt_LargeData(t *testing.T) {
	key, _ := aes.GenerateAESKey(256)
	// Test with 1MB of data
	largeText := strings.Repeat("A", 1024*1024)

	ciphertext, err := aes.Encrypt(key, largeText)
	require.NoError(t, err)

	decrypted, err := aes.Decrypt(key, ciphertext)
	require.NoError(t, err)
	require.Equal(t, largeText, decrypted)
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1, _ := aes.GenerateAESKey(256)
	key2, _ := aes.GenerateAESKey(256)

	plaintext := "secret message"
	ciphertext, err := aes.Encrypt(key1, plaintext)
	require.NoError(t, err)

	// Try to decrypt with wrong key
	_, err = aes.Decrypt(key2, ciphertext)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt")
}

func TestDecrypt_InvalidKeyLength(t *testing.T) {
	// Create invalid key lengths for AES (valid AES key lengths are 16, 24, 32 bytes)
	invalidKeys := [][]byte{
		make([]byte, 10), // 10 bytes
		make([]byte, 15), // 15 bytes
		make([]byte, 17), // 17 bytes
		make([]byte, 20), // 20 bytes
		make([]byte, 1),  // 1 byte
		make([]byte, 33), // 33 bytes
	}

	for _, keyBytes := range invalidKeys {
		invalidKey := base64.StdEncoding.EncodeToString(keyBytes)
		_, err := aes.Encrypt(invalidKey, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create AES cipher")
	}
}

func TestEncrypt_CiphertextRandomness(t *testing.T) {
	key, _ := aes.GenerateAESKey(256)
	plaintext := "same message"

	// Encrypt same message multiple times, should get different ciphertexts
	ciphertext1, err := aes.Encrypt(key, plaintext)
	require.NoError(t, err)

	ciphertext2, err := aes.Encrypt(key, plaintext)
	require.NoError(t, err)

	require.NotEqual(t, ciphertext1, ciphertext2, "Same plaintext should produce different ciphertexts due to random nonce")

	// Both should decrypt to same plaintext
	decrypted1, err := aes.Decrypt(key, ciphertext1)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted1)

	decrypted2, err := aes.Decrypt(key, ciphertext2)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted2)
}
