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

func TestGenerateHashedPassword_Argon2i(t *testing.T) {
	cfg := argon.Config{
		Type:        "argon2i",
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := argon.GenerateHashedPassword("test-password", cfg)
	require.NoError(t, err)
	require.Contains(t, hash, "$argon2i$")

	ok, err := argon.CompareHashAndPassword(hash, "test-password")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestGenerateHashedPassword_UnsupportedType(t *testing.T) {
	cfg := argon.Config{
		Type:        "argon2d",
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	_, err := argon.GenerateHashedPassword("password", cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported Argon2 type")
}

func TestGenerateHashedPassword_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		cfg       argon.Config
		shouldErr bool
	}{
		{
			name: "minimal valid config",
			cfg: argon.Config{
				Type:        "argon2id",
				Memory:      8,
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  8,
				KeyLength:   8,
			},
			shouldErr: false,
		},
		{
			name: "very small salt",
			cfg: argon.Config{
				Type:        "argon2id",
				Memory:      1024,
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  1,
				KeyLength:   32,
			},
			shouldErr: false,
		},
		{
			name: "very small key",
			cfg: argon.Config{
				Type:        "argon2id",
				Memory:      1024,
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  16,
				KeyLength:   1,
			},
			shouldErr: false,
		},
		{
			name: "large key and salt",
			cfg: argon.Config{
				Type:        "argon2id",
				Memory:      1024,
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  256,
				KeyLength:   256,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := argon.GenerateHashedPassword("password", tt.cfg)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, hash)
				
				// Verify the hash can be used for comparison
				ok, err := argon.CompareHashAndPassword(hash, "password")
				require.NoError(t, err)
				require.True(t, ok)
			}
		})
	}
}

func TestCompareHashAndPassword_HashFormatErrors(t *testing.T) {
	tests := []struct {
		name string
		hash string
		err  error
	}{
		{
			name: "too few parts",
			hash: "$argon2id$v=19$m=65536",
			err:  argon.ErrInvalidHashFormat,
		},
		{
			name: "too many parts",
			hash: "$argon2id$v=19$m=65536,t=2,p=2$salt$hash$extra",
			err:  argon.ErrInvalidHashFormat,
		},
		{
			name: "empty hash",
			hash: "",
			err:  argon.ErrInvalidHashFormat,
		},
		{
			name: "no dollar signs",
			hash: "argon2id_v=19_m=65536,t=2,p=2_salt_hash",
			err:  argon.ErrInvalidHashFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := argon.CompareHashAndPassword(tt.hash, "password")
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestCompareHashAndPassword_ParameterParsingErrors(t *testing.T) {
	tests := []struct {
		name string
		hash string
		errMsg string
	}{
		{
			name: "invalid version format",
			hash: "$argon2id$version=19$m=65536,t=2,p=2$dGVzdA$dGVzdA",
			errMsg: "version mismatch",
		},
		{
			name: "non-numeric version",
			hash: "$argon2id$v=abc$m=65536,t=2,p=2$dGVzdA$dGVzdA",
			errMsg: "version mismatch",
		},
		{
			name: "too few parameters",
			hash: "$argon2id$v=19$m=65536,t=2$dGVzdA$dGVzdA",
			errMsg: "invalid password hash format",
		},
		{
			name: "too many parameters",
			hash: "$argon2id$v=19$m=65536,t=2,p=2,x=1$dGVzdA$dGVzdA",
			errMsg: "invalid password hash format",
		},
		{
			name: "invalid memory parameter",
			hash: "$argon2id$v=19$m=abc,t=2,p=2$dGVzdA$dGVzdA",
			errMsg: "invalid memory parameter",
		},
		{
			name: "invalid iterations parameter",
			hash: "$argon2id$v=19$m=65536,t=abc,p=2$dGVzdA$dGVzdA",
			errMsg: "invalid iterations parameter",
		},
		{
			name: "invalid parallelism parameter",
			hash: "$argon2id$v=19$m=65536,t=2,p=abc$dGVzdA$dGVzdA",
			errMsg: "invalid parallelism parameter",
		},
		{
			name: "negative memory",
			hash: "$argon2id$v=19$m=-1000,t=2,p=2$dGVzdA$dGVzdA",
			errMsg: "invalid memory size",
		},
		{
			name: "negative iterations",
			hash: "$argon2id$v=19$m=65536,t=-5,p=2$dGVzdA$dGVzdA",
			errMsg: "invalid iterations count",
		},
		{
			name: "negative parallelism",
			hash: "$argon2id$v=19$m=65536,t=2,p=-3$dGVzdA$dGVzdA",
			errMsg: "invalid parallelism count",
		},
		{
			name: "parallelism too large",
			hash: "$argon2id$v=19$m=65536,t=2,p=300$dGVzdA$dGVzdA",
			errMsg: "invalid parallelism count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := argon.CompareHashAndPassword(tt.hash, "password")
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestCompareHashAndPassword_EmptyPassword(t *testing.T) {
	cfg := argon.DefaultConfig
	hash, err := argon.GenerateHashedPassword("nonempty", cfg)
	require.NoError(t, err)

	// Empty password should fail comparison
	ok, err := argon.CompareHashAndPassword(hash, "")
	require.ErrorIs(t, err, argon.ErrInvalidPassword)
	require.False(t, ok)
}

func TestDefaultConfig(t *testing.T) {
	cfg := argon.DefaultConfig
	require.Equal(t, "argon2id", cfg.Type)
	require.Equal(t, uint32(128*1024), cfg.Memory)
	require.Equal(t, uint32(4), cfg.Iterations)
	require.Greater(t, cfg.Parallelism, uint8(0))
	require.Equal(t, uint32(32), cfg.SaltLength)
	require.Equal(t, uint32(32), cfg.KeyLength)
}

func TestGenerateHashedPassword_ConsistentFormat(t *testing.T) {
	cfg := argon.DefaultConfig
	hash, err := argon.GenerateHashedPassword("test", cfg)
	require.NoError(t, err)

	// Check hash format structure
	parts := strings.Split(hash, "$")
	require.Len(t, parts, 6)
	require.Equal(t, "", parts[0]) // Leading empty part from first $
	require.Equal(t, "argon2id", parts[1])
	require.True(t, strings.HasPrefix(parts[2], "v="))
	require.Contains(t, parts[3], "m=")
	require.Contains(t, parts[3], "t=")
	require.Contains(t, parts[3], "p=")
	require.NotEmpty(t, parts[4]) // salt
	require.NotEmpty(t, parts[5]) // hash
}

func TestCompareHashAndPassword_DifferentPasswords(t *testing.T) {
	cfg := argon.DefaultConfig
	hash, err := argon.GenerateHashedPassword("original", cfg)
	require.NoError(t, err)

	testCases := []string{
		"different",
		"Original", // case sensitive
		"original ", // trailing space
		" original", // leading space
		"",
	}

	for _, wrongPassword := range testCases {
		ok, err := argon.CompareHashAndPassword(hash, wrongPassword)
		require.ErrorIs(t, err, argon.ErrInvalidPassword)
		require.False(t, ok)
	}
}

func TestGenerateHashedPassword_UniqueHashes(t *testing.T) {
	cfg := argon.DefaultConfig
	password := "same-password"

	hash1, err := argon.GenerateHashedPassword(password, cfg)
	require.NoError(t, err)

	hash2, err := argon.GenerateHashedPassword(password, cfg)
	require.NoError(t, err)

	// Should be different due to random salt
	require.NotEqual(t, hash1, hash2)

	// But both should verify correctly
	ok1, err := argon.CompareHashAndPassword(hash1, password)
	require.NoError(t, err)
	require.True(t, ok1)

	ok2, err := argon.CompareHashAndPassword(hash2, password)
	require.NoError(t, err)
	require.True(t, ok2)
}
