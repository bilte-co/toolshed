// Package aes provides AES-GCM encryption and decryption functionality with base64 encoding.
// It supports 128, 192, and 256-bit keys and uses authenticated encryption for security.
//
// Example usage:
//
//	// Generate a key
//	key, err := aes.GenerateAESKey(256)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Encrypt data
//	ciphertext, err := aes.Encrypt(key, "secret message")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Decrypt data
//	plaintext, err := aes.Decrypt(key, ciphertext)
//	if err != nil {
//		log.Fatal(err)
//	}
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// GenerateAESKey creates a base64-encoded AES key.
// Valid bit lengths: 128, 192, or 256.
func GenerateAESKey(bits int) (string, error) {
	if bits != 128 && bits != 192 && bits != 256 {
		return "", fmt.Errorf("invalid AES key length: %d (must be 128, 192, or 256)", bits)
	}

	key := make([]byte, bits/8)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// Encrypt encrypts the plaintext using AES-GCM and returns a base64-encoded ciphertext.
func Encrypt(b64Key string, plaintext string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return "", fmt.Errorf("invalid base64 key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AES-GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext using the provided base64 key.
func Decrypt(b64Key string, b64Ciphertext string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return "", fmt.Errorf("invalid base64 key: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(b64Ciphertext)
	if err != nil {
		return "", fmt.Errorf("invalid base64 ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AES-GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short: missing nonce")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(plaintext), nil
}
