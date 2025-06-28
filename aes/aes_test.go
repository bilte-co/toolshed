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
