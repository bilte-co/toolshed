// Package argon provides secure Argon2 password hashing functionality.
// It supports both Argon2id (recommended) and Argon2i variants with configurable parameters.
// 
// Argon2 is a memory-hard key derivation function that won the Password Hashing Competition
// and is resistant to both time-memory trade-off attacks and side-channel attacks.
//
// Example usage:
//
//	// Hash a password with default configuration
//	hash, err := argon.GenerateHashedPassword("mypassword", argon.DefaultConfig)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Verify a password
//	valid, err := argon.CompareHashAndPassword(hash, "mypassword")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if valid {
//		fmt.Println("Password is correct")
//	}
package argon

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Config defines tunable parameters for Argon2 password hashing.
type Config struct {
	Type        string // "argon2id" or "argon2i"
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultConfig provides sane, secure default values.
var DefaultConfig = Config{
	Type:        "argon2id",
	Memory:      128 * 1024, // 128 MiB
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  32,
	KeyLength:   32,
}

// Exported error types for use in conditional handling
var (
	ErrInvalidHashFormat = errors.New("invalid password hash format")
	ErrVersionMismatch   = errors.New("argon2 version mismatch")
	ErrInvalidPassword   = errors.New("password does not match")
)

// GenerateHashedPassword hashes the given password using the provided Argon2 configuration.
func GenerateHashedPassword(password string, cfg Config) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
	}

	salt, err := generateRandomBytes(int(cfg.SaltLength))
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hashBytes, err := hashPassword(cfg, salt, password)
	if err != nil {
		return "", err
	}

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hashBytes)

	hashString := fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		cfg.Type, argon2.Version, cfg.Memory, cfg.Iterations, cfg.Parallelism,
		saltBase64, hashBase64)

	return hashString, nil
}

// CompareHashAndPassword verifies a password against a stored hash.
func CompareHashAndPassword(storedHash, password string) (bool, error) {
	cfg, salt, expectedHash, err := parseHash(storedHash)
	if err != nil {
		return false, err
	}

	calculatedHash, err := hashPassword(cfg, salt, password)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare(expectedHash, calculatedHash) != 1 {
		return false, ErrInvalidPassword
	}
	return true, nil
}

// generateRandomBytes securely generates a random byte slice of given length.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// hashPassword generates the Argon2 hash based on the config.
func hashPassword(cfg Config, salt []byte, password string) ([]byte, error) {
	switch cfg.Type {
	case "argon2id":
		return argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength), nil
	case "argon2i":
		return argon2.Key([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength), nil
	default:
		return nil, fmt.Errorf("unsupported Argon2 type: %q", cfg.Type)
	}
}

// parseHash parses the Argon2 hash string into parameters, salt, and hash.
func parseHash(hash string) (Config, []byte, []byte, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return Config{}, nil, nil, ErrInvalidHashFormat
	}

	argonType := parts[1]
	versionStr := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionStr)
	if err != nil || version != argon2.Version {
		return Config{}, nil, nil, ErrVersionMismatch
	}

	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return Config{}, nil, nil, ErrInvalidHashFormat
	}

	memory, err := strconv.Atoi(strings.TrimPrefix(params[0], "m="))
	if err != nil {
		return Config{}, nil, nil, errors.New("invalid memory parameter")
	}

	iterations, err := strconv.Atoi(strings.TrimPrefix(params[1], "t="))
	if err != nil {
		return Config{}, nil, nil, errors.New("invalid iterations parameter")
	}

	parallelism, err := strconv.Atoi(strings.TrimPrefix(params[2], "p="))
	if err != nil {
		return Config{}, nil, nil, errors.New("invalid parallelism parameter")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Config{}, nil, nil, errors.New("invalid base64 salt")
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Config{}, nil, nil, errors.New("invalid base64 hash")
	}

	cfg := Config{
		Type:        argonType,
		Memory:      uint32(memory),
		Iterations:  uint32(iterations),
		Parallelism: uint8(parallelism),
		SaltLength:  uint32(len(salt)),
		KeyLength:   uint32(len(hashBytes)),
	}

	return cfg, salt, hashBytes, nil
}
