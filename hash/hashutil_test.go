package hash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test vectors from known sources
var testVectors = map[string]map[string]string{
	"sha256": {
		"":      "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"hello": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		"world": "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7",
	},
	"sha512": {
		"": "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		"hello": "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
	},
	"md5": {
		"":      "d41d8cd98f00b204e9800998ecf8427e",
		"hello": "5d41402abc4b2a76b9719d911017c592",
	},
	"sha1": {
		"":      "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		"hello": "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
	},
}

func TestHashString(t *testing.T) {
	for algo, vectors := range testVectors {
		t.Run(algo, func(t *testing.T) {
			for input, expected := range vectors {
				result, err := HashString(input, algo)
				require.NoError(t, err, "HashString should not return error for %s", algo)
				
				actual := hex.EncodeToString(result)
				assert.Equal(t, expected, actual, "Hash mismatch for algorithm %s, input %q", algo, input)
			}
		})
	}
}

func TestHashBytes(t *testing.T) {
	for algo, vectors := range testVectors {
		t.Run(algo, func(t *testing.T) {
			for input, expected := range vectors {
				result, err := HashBytes([]byte(input), algo)
				require.NoError(t, err, "HashBytes should not return error for %s", algo)
				
				actual := hex.EncodeToString(result)
				assert.Equal(t, expected, actual, "Hash mismatch for algorithm %s, input %q", algo, input)
			}
		})
	}
}

func TestHashReader(t *testing.T) {
	for algo, vectors := range testVectors {
		t.Run(algo, func(t *testing.T) {
			for input, expected := range vectors {
				reader := strings.NewReader(input)
				result, err := HashReader(reader, algo)
				require.NoError(t, err, "HashReader should not return error for %s", algo)
				
				actual := hex.EncodeToString(result)
				assert.Equal(t, expected, actual, "Hash mismatch for algorithm %s, input %q", algo, input)
			}
		})
	}
}

func TestHashFile(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "hello world"
	
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Failed to create test file")
	
	// Test hashing the file
	result, err := HashFile(testFile, "sha256")
	require.NoError(t, err, "HashFile should not return error")
	
	// Compare with known hash of "hello world"
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	actual := hex.EncodeToString(result)
	assert.Equal(t, expected, actual, "File hash mismatch")
}

func TestHashDir(t *testing.T) {
	// Create temporary test directory structure
	tmpDir := t.TempDir()
	
	// Create test files
	files := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"subdir/file3.txt": "content3",
	}
	
	for filePath, content := range files {
		fullPath := filepath.Join(tmpDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err, "Failed to create directory")
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err, "Failed to create test file")
	}
	
	// Test non-recursive hashing (should only include files in root)
	hash1, err := HashDir(tmpDir, "sha256", false)
	require.NoError(t, err, "HashDir should not return error")
	
	// Test recursive hashing (should include all files)
	hash2, err := HashDir(tmpDir, "sha256", true)
	require.NoError(t, err, "HashDir should not return error")
	
	// Hashes should be different (recursive includes more files)
	assert.NotEqual(t, hash1, hash2, "Recursive and non-recursive hashes should differ")
	
	// Test deterministic behavior
	hash3, err := HashDir(tmpDir, "sha256", true)
	require.NoError(t, err, "HashDir should not return error")
	assert.Equal(t, hash2, hash3, "HashDir should be deterministic")
}

func TestUnsupportedAlgorithm(t *testing.T) {
	_, err := HashString("test", "unsupported")
	assert.Error(t, err, "Should return error for unsupported algorithm")
	assert.Contains(t, err.Error(), "unsupported hash algorithm")
}

func TestHashStringWithOptions(t *testing.T) {
	input := "test"
	
	// Test hex format
	result, err := HashStringWithOptions(input, "sha256", Options{Format: FormatHex})
	require.NoError(t, err, "Should not return error")
	assert.IsType(t, "", result, "Result should be string for hex format")
	
	// Test base64 format
	result, err = HashStringWithOptions(input, "sha256", Options{Format: FormatBase64})
	require.NoError(t, err, "Should not return error")
	assert.IsType(t, "", result, "Result should be string for base64 format")
	
	// Test raw format
	result, err = HashStringWithOptions(input, "sha256", Options{Format: FormatRaw})
	require.NoError(t, err, "Should not return error")
	assert.IsType(t, []byte{}, result, "Result should be []byte for raw format")
	
	// Test with prefix
	result, err = HashStringWithOptions(input, "sha256", Options{Format: FormatHex, Prefix: true})
	require.NoError(t, err, "Should not return error")
	resultStr := result.(string)
	assert.True(t, strings.HasPrefix(resultStr, "sha256:"), "Result should have algorithm prefix")
}

func TestNewHasher(t *testing.T) {
	hasher, err := NewHasher("sha256")
	require.NoError(t, err, "Should create hasher successfully")
	
	assert.Equal(t, "sha256", hasher.Algorithm(), "Algorithm should match")
	
	// Test incremental hashing
	hasher.Write([]byte("hello"))
	hasher.Write([]byte(" "))
	hasher.Write([]byte("world"))
	
	result := hasher.SumHex()
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	assert.Equal(t, expected, result, "Incremental hashing should work correctly")
	
	// Test base64 output
	base64Result := hasher.SumBase64()
	assert.NotEmpty(t, base64Result, "Base64 result should not be empty")
}

func TestRegisterHasher(t *testing.T) {
	// Register a custom hasher (using existing sha256 for simplicity)
	RegisterHasher("custom", func() hash.Hash {
		return sha256.New()
	})
	
	// Test using the custom hasher
	result, err := HashString("test", "custom")
	require.NoError(t, err, "Should work with registered custom hasher")
	
	// Should produce same result as sha256
	expected, err := HashString("test", "sha256")
	require.NoError(t, err, "SHA256 should work")
	
	assert.Equal(t, expected, result, "Custom hasher should produce same result as sha256")
}

func TestBLAKE2b(t *testing.T) {
	// Test BLAKE2b algorithm
	result, err := HashString("test", "blake2b")
	require.NoError(t, err, "BLAKE2b should work")
	assert.NotEmpty(t, result, "BLAKE2b result should not be empty")
	assert.Len(t, result, 32, "BLAKE2b-256 should produce 32-byte hash")
}

func TestSecurityWarnings(t *testing.T) {
	// This test ensures that security warnings are triggered for insecure algorithms
	// We can't easily test log output, but we can ensure the functions still work
	
	_, err := HashString("test", "md5")
	assert.NoError(t, err, "MD5 should work but with warning")
	
	_, err = HashString("test", "sha1")
	assert.NoError(t, err, "SHA1 should work but with warning")
}

func BenchmarkHashString(b *testing.B) {
	data := make([]byte, 1024) // 1KB of data
	rand.Read(data)
	input := string(data)
	
	algorithms := []string{"sha256", "sha512", "blake2b"}
	
	for _, algo := range algorithms {
		b.Run(algo, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := HashString(input, algo)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkHashReader(b *testing.B) {
	data := make([]byte, 1024*1024) // 1MB of data
	rand.Read(data)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		_, err := HashReader(reader, "sha256")
		if err != nil {
			b.Fatal(err)
		}
	}
}
