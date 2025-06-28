package hash

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHMAC(t *testing.T) {
	data := []byte("hello world")
	key := []byte("secret key")

	// Test HMAC with different algorithms
	algorithms := []string{"sha256", "sha512", "sha1"}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			result, err := HMAC(data, key, algo)
			require.NoError(t, err, "HMAC should not return error")
			assert.NotEmpty(t, result, "HMAC result should not be empty")

			// Test consistency
			result2, err := HMAC(data, key, algo)
			require.NoError(t, err, "HMAC should not return error")
			assert.Equal(t, result, result2, "HMAC should be deterministic")

			// Test with different key
			result3, err := HMAC(data, []byte("different key"), algo)
			require.NoError(t, err, "HMAC should not return error")
			assert.NotEqual(t, result, result3, "HMAC should differ with different keys")
		})
	}
}

func TestHMACWithOptions(t *testing.T) {
	data := []byte("test data")
	key := []byte("test key")

	// Test with hex format
	result, err := HMACWithOptions(data, key, "sha256", Options{Format: FormatHex})
	require.NoError(t, err, "HMAC with options should not return error")
	assert.IsType(t, "", result, "Result should be string for hex format")

	// Test with prefix
	result, err = HMACWithOptions(data, key, "sha256", Options{Format: FormatHex, Prefix: true})
	require.NoError(t, err, "HMAC with options should not return error")
	resultStr := result.(string)
	assert.Contains(t, resultStr, "hmac-sha256:", "Result should have HMAC prefix")
}

func TestEqualConstantTime(t *testing.T) {
	// Test equal byte slices
	a := []byte("hello")
	b := []byte("hello")
	assert.True(t, EqualConstantTime(a, b), "Equal byte slices should return true")

	// Test different byte slices
	c := []byte("world")
	assert.False(t, EqualConstantTime(a, c), "Different byte slices should return false")

	// Test different lengths
	d := []byte("hello world")
	assert.False(t, EqualConstantTime(a, d), "Different length slices should return false")

	// Test empty slices
	assert.True(t, EqualConstantTime([]byte{}, []byte{}), "Empty slices should be equal")
	assert.False(t, EqualConstantTime(a, []byte{}), "Non-empty and empty should not be equal")
}

func TestPBKDF2Hash(t *testing.T) {
	password := []byte("testpassword")
	opts := &PasswordHashingOptions{
		PBKDF2Iterations: 1000,
		PBKDF2KeyLength:  32,
		PBKDF2SaltLength: 16,
		PBKDF2Algorithm:  "sha256",
	}

	key, salt, err := PBKDF2Hash(password, opts)
	require.NoError(t, err, "PBKDF2Hash should not return error")
	assert.Len(t, key, 32, "Key should be 32 bytes")
	assert.Len(t, salt, 16, "Salt should be 16 bytes")

	// Test with same password and salt should produce same key
	key2, _, err := PBKDF2HashWithSalt(password, salt, opts)
	require.NoError(t, err, "PBKDF2HashWithSalt should not return error")
	assert.Equal(t, key, key2, "Same password and salt should produce same key")

	// Test verification
	assert.True(t, VerifyPBKDF2(password, salt, key, opts), "Verification should succeed")
	assert.False(t, VerifyPBKDF2([]byte("wrongpassword"), salt, key, opts), "Wrong password should fail verification")
}

func TestPBKDF2DefaultOptions(t *testing.T) {
	password := []byte("testpassword")

	key, salt, err := PBKDF2Hash(password, nil)
	require.NoError(t, err, "PBKDF2Hash with nil options should work")
	assert.NotEmpty(t, key, "Key should not be empty")
	assert.NotEmpty(t, salt, "Salt should not be empty")
}

func TestScryptHash(t *testing.T) {
	password := []byte("testpassword")
	opts := &PasswordHashingOptions{
		ScryptN:       1024, // Lower for testing
		ScryptR:       8,
		ScryptP:       1,
		ScryptKeyLen:  32,
		ScryptSaltLen: 16,
	}

	key, salt, err := ScryptHash(password, opts)
	require.NoError(t, err, "ScryptHash should not return error")
	assert.Len(t, key, 32, "Key should be 32 bytes")
	assert.Len(t, salt, 16, "Salt should be 16 bytes")

	// Test with same password and salt should produce same key
	key2, _, err := ScryptHashWithSalt(password, salt, opts)
	require.NoError(t, err, "ScryptHashWithSalt should not return error")
	assert.Equal(t, key, key2, "Same password and salt should produce same key")

	// Test verification
	assert.True(t, VerifyScrypt(password, salt, key, opts), "Verification should succeed")
	assert.False(t, VerifyScrypt([]byte("wrongpassword"), salt, key, opts), "Wrong password should fail verification")
}

func TestScryptDefaultOptions(t *testing.T) {
	password := []byte("testpassword")

	// Use lower N value for testing to avoid long test times
	opts := DefaultPasswordOptions
	opts.ScryptN = 1024

	key, salt, err := ScryptHash(password, &opts)
	require.NoError(t, err, "ScryptHash with default options should work")
	assert.NotEmpty(t, key, "Key should not be empty")
	assert.NotEmpty(t, salt, "Salt should not be empty")
}

func TestBcryptHash(t *testing.T) {
	password := []byte("testpassword")
	opts := &PasswordHashingOptions{
		BcryptCost: 8, // Lower cost for testing
	}

	hash, err := BcryptHash(password, opts)
	require.NoError(t, err, "BcryptHash should not return error")
	assert.NotEmpty(t, hash, "Hash should not be empty")

	// Test verification
	assert.True(t, VerifyBcrypt(password, hash), "Verification should succeed")
	assert.False(t, VerifyBcrypt([]byte("wrongpassword"), hash), "Wrong password should fail verification")

	// Test that different calls produce different hashes (due to random salt)
	hash2, err := BcryptHash(password, opts)
	require.NoError(t, err, "BcryptHash should not return error")
	assert.NotEqual(t, hash, hash2, "Different bcrypt calls should produce different hashes")

	// But both should verify correctly
	assert.True(t, VerifyBcrypt(password, hash2), "Second hash should also verify")
}

func TestBcryptDefaultOptions(t *testing.T) {
	password := []byte("testpassword")

	hash, err := BcryptHash(password, nil)
	require.NoError(t, err, "BcryptHash with nil options should work")
	assert.NotEmpty(t, hash, "Hash should not be empty")
	assert.True(t, VerifyBcrypt(password, hash), "Verification should succeed")
}

func TestGenerateSalt(t *testing.T) {
	// Test different salt lengths
	lengths := []int{16, 32, 64}

	for _, length := range lengths {
		salt, err := GenerateSalt(length)
		require.NoError(t, err, "GenerateSalt should not return error")
		assert.Len(t, salt, length, "Salt should have correct length")

		// Test that multiple calls produce different salts
		salt2, err := GenerateSalt(length)
		require.NoError(t, err, "GenerateSalt should not return error")
		assert.NotEqual(t, salt, salt2, "Different calls should produce different salts")
	}
}

func TestPasswordHashingCrossPlatform(t *testing.T) {
	// Test that our implementations work consistently
	password := []byte("test123")

	// PBKDF2 test
	opts := &PasswordHashingOptions{
		PBKDF2Iterations: 1000,
		PBKDF2KeyLength:  32,
		PBKDF2SaltLength: 16,
		PBKDF2Algorithm:  "sha256",
	}

	salt := []byte("1234567890123456") // Fixed salt for reproducibility
	key, _, err := PBKDF2HashWithSalt(password, salt, opts)
	require.NoError(t, err, "PBKDF2HashWithSalt should not return error")

	// The key should be deterministic with fixed salt
	key2, _, err := PBKDF2HashWithSalt(password, salt, opts)
	require.NoError(t, err, "PBKDF2HashWithSalt should not return error")
	assert.Equal(t, key, key2, "PBKDF2 should be deterministic with same inputs")
}

func BenchmarkHMAC(b *testing.B) {
	data := make([]byte, 1024)
	key := make([]byte, 32)
	rand.Read(data)
	rand.Read(key)

	for b.Loop() {
		_, err := HMAC(data, key, "sha256")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	password := []byte("testpassword")
	salt := make([]byte, 16)
	rand.Read(salt)

	opts := &PasswordHashingOptions{
		PBKDF2Iterations: 10000,
		PBKDF2KeyLength:  32,
		PBKDF2Algorithm:  "sha256",
	}

	for b.Loop() {
		_, _, err := PBKDF2HashWithSalt(password, salt, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBcrypt(b *testing.B) {
	password := []byte("testpassword")
	opts := &PasswordHashingOptions{
		BcryptCost: 10,
	}

	for b.Loop() {
		_, err := BcryptHash(password, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
