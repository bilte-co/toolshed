package argon_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bilte-co/toolshed/argon"
)

func TestGenerateAndCompare_ValidPassword(t *testing.T) {
	cfg := argon.DefaultConfig
	password := "correct horse battery staple"

	hash, err := argon.GenerateHashedPassword(password, cfg)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	ok, err := argon.CompareHashAndPassword(hash, password)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestCompareHashAndPassword_WrongPassword(t *testing.T) {
	cfg := argon.DefaultConfig
	hash, err := argon.GenerateHashedPassword("secret123", cfg)
	require.NoError(t, err)

	ok, err := argon.CompareHashAndPassword(hash, "wrong-password")
	require.ErrorIs(t, err, argon.ErrInvalidPassword)
	require.False(t, ok)
}

func TestGenerateHashedPassword_EmptyPassword(t *testing.T) {
	_, err := argon.GenerateHashedPassword("", argon.DefaultConfig)
	require.Error(t, err)
	require.Contains(t, err.Error(), "password cannot be empty")
}

func TestCompareHashAndPassword_MalformedHash(t *testing.T) {
	_, err := argon.CompareHashAndPassword("bad$hash$format", "password")
	require.ErrorIs(t, err, argon.ErrInvalidHashFormat)
}

func TestCompareHashAndPassword_UnsupportedType(t *testing.T) {
	hash := "$argon3id$v=19$m=65536,t=2,p=2$" + strings.Repeat("a", 32) + "$" + strings.Repeat("b", 32)
	_, err := argon.CompareHashAndPassword(hash, "any-password")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported Argon2 type")
}

func TestCompareHashAndPassword_InvalidVersion(t *testing.T) {
	hash := "$argon2id$v=999$m=65536,t=2,p=2$" + strings.Repeat("a", 32) + "$" + strings.Repeat("b", 32)
	_, err := argon.CompareHashAndPassword(hash, "any-password")
	require.ErrorIs(t, err, argon.ErrVersionMismatch)
}

func TestCompareHashAndPassword_BadBase64Salt(t *testing.T) {
	hash := "$argon2id$v=19$m=65536,t=2,p=2$%%%invalidbase64%%%$" + strings.Repeat("b", 32)
	_, err := argon.CompareHashAndPassword(hash, "password")
	require.ErrorContains(t, err, "invalid base64 salt")
}

func TestCompareHashAndPassword_BadBase64Hash(t *testing.T) {
	hash := "$argon2id$v=19$m=65536,t=2,p=2$" + strings.Repeat("a", 32) + "$%%%bad%%%"
	_, err := argon.CompareHashAndPassword(hash, "password")
	require.ErrorContains(t, err, "invalid base64 hash")
}

func TestGenerateHashedPassword_CustomConfig(t *testing.T) {
	cfg := argon.Config{
		Type:        "argon2id",
		Memory:      64 * 1024,
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   24,
	}

	hash, err := argon.GenerateHashedPassword("custom-pass", cfg)
	require.NoError(t, err)

	ok, err := argon.CompareHashAndPassword(hash, "custom-pass")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestParsedConfigFromHashMatches(t *testing.T) {
	cfg := argon.Config{
		Type:        "argon2id",
		Memory:      96 * 1024,
		Iterations:  3,
		Parallelism: 3,
		SaltLength:  20,
		KeyLength:   32,
	}

	hash, err := argon.GenerateHashedPassword("testpass", cfg)
	require.NoError(t, err)

	ok, err := argon.CompareHashAndPassword(hash, "testpass")
	require.NoError(t, err)
	require.True(t, ok)
}
