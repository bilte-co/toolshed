// Package base64 provides convenient wrapper functions for base64 encoding and decoding.
// All functions use the standard base64 encoding (RFC 4648) which is URL-safe and
// widely compatible across systems.
//
// Example usage:
//
//	// Encode a string to base64
//	encoded := base64.EncodeString("Hello, World!")
//	fmt.Println(encoded) // SGVsbG8sIFdvcmxkIQ==
//
//	// Decode a base64 string
//	decoded, err := base64.DecodeToString(encoded)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(decoded) // Hello, World!
package base64

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Encode encodes the given data to base64 string
func Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// EncodeString encodes the given string to base64 string
func EncodeString(s string) string {
	return Encode([]byte(s))
}

// Decode decodes the given base64 string to bytes
func Decode(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return nil, fmt.Errorf("invalid base64 input: %v", err)
	}
	return decoded, nil
}

// DecodeToString decodes the given base64 string to string
func DecodeToString(encoded string) (string, error) {
	decoded, err := Decode(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
