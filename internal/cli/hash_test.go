package cli_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/bilte-co/toolshed/hash"
	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestHashStringCmd_BasicOperation(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		algo     string
		format   string
		expected string // known hash for verification
	}{
		{
			name:     "empty string sha256",
			text:     "",
			algo:     "sha256",
			format:   "hex",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello sha256",
			text:     "hello",
			algo:     "sha256",
			format:   "hex",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "hello md5",
			text:     "hello", 
			algo:     "md5",
			format:   "hex",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HashStringCmd{
				Text:   tt.text,
				Algo:   tt.algo,
				Format: tt.format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHashStringCmd_AllAlgorithms(t *testing.T) {
	algorithms := []string{"md5", "sha1", "sha256", "sha512", "blake2b"}
	testText := "test message"

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			cmd := &cli.HashStringCmd{
				Text: testText,
				Algo: algo,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHashStringCmd_AllFormats(t *testing.T) {
	formats := []string{"hex", "base64", "raw"}
	testText := "format test"

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			cmd := &cli.HashStringCmd{
				Text:   testText,
				Format: format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHashStringCmd_WithPrefix(t *testing.T) {
	cmd := &cli.HashStringCmd{
		Text:   "test",
		Algo:   "sha256",
		Format: "hex",
		Prefix: true,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashStringCmd_InvalidAlgorithm(t *testing.T) {
	cmd := &cli.HashStringCmd{
		Text: "test",
		Algo: "invalid-algo",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.Error(t, err)
}

func TestHashFileCmd_RegularFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "file content for hashing"

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	cmd := &cli.HashFileCmd{
		Path: testFile,
		Algo: "sha256",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashFileCmd_StdinInput(t *testing.T) {
	cmd := &cli.HashFileCmd{
		Path: "-",
		Algo: "sha256",
	}
	ctx := newTestContext()

	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	testContent := "stdin content"

	go func() {
		defer w.Close()
		w.Write([]byte(testContent))
	}()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashFileCmd_NonexistentFile(t *testing.T) {
	cmd := &cli.HashFileCmd{
		Path: "/nonexistent/file.txt",
		Algo: "sha256",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.Error(t, err)
}

func TestHashFileCmd_WithAllOptions(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "comprehensive test content"

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	cmd := &cli.HashFileCmd{
		Path:   testFile,
		Algo:   "sha512",
		Format: "base64",
		Prefix: true,
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashDirCmd_BasicDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create some files in the directory
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	
	err := os.WriteFile(file1, []byte("content1"), 0o644)
	require.NoError(t, err)
	
	err = os.WriteFile(file2, []byte("content2"), 0o644)
	require.NoError(t, err)

	cmd := &cli.HashDirCmd{
		Path: tmpDir,
		Algo: "sha256",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashDirCmd_RecursiveVsNonRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	
	err := os.Mkdir(subDir, 0o755)
	require.NoError(t, err)
	
	// Create files in both directories
	err = os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0o644)
	require.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(subDir, "sub.txt"), []byte("sub"), 0o644)
	require.NoError(t, err)

	// Test recursive (default)
	cmdRecursive := &cli.HashDirCmd{
		Path:      tmpDir,
		Algo:      "sha256",
		Recursive: true,
	}
	ctx := newTestContext()

	err = cmdRecursive.Run(ctx)
	require.NoError(t, err)

	// Test non-recursive
	cmdNonRecursive := &cli.HashDirCmd{
		Path:      tmpDir,
		Algo:      "sha256",
		Recursive: false,
	}

	err = cmdNonRecursive.Run(ctx)
	require.NoError(t, err)
}

func TestHashDirCmd_NonexistentDirectory(t *testing.T) {
	cmd := &cli.HashDirCmd{
		Path: "/nonexistent/directory",
		Algo: "sha256",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.Error(t, err)
}

func TestHMACCmd_BasicOperation(t *testing.T) {
	cmd := &cli.HMACCmd{
		Text: "test message",
		Key:  "secret-key",
		Algo: "sha256",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHMACCmd_AllAlgorithms(t *testing.T) {
	algorithms := []string{"md5", "sha1", "sha256", "sha512", "blake2b"}
	testText := "hmac test message"
	testKey := "hmac-test-key"

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			cmd := &cli.HMACCmd{
				Text: testText,
				Key:  testKey,
				Algo: algo,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHMACCmd_AllFormats(t *testing.T) {
	formats := []string{"hex", "base64", "raw"}
	testText := "hmac format test"
	testKey := "format-key"

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			cmd := &cli.HMACCmd{
				Text:   testText,
				Key:    testKey,
				Format: format,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHMACCmd_WithPrefix(t *testing.T) {
	cmd := &cli.HMACCmd{
		Text:   "prefix test",
		Key:    "prefix-key",
		Algo:   "sha256",
		Format: "hex",
		Prefix: true,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestValidateCmd_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "validate.txt")
	testContent := "content to validate"

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	// Calculate expected hash
	expectedHash, err := hash.HashString(testContent, "sha256")
	require.NoError(t, err)
	expectedHex := hex.EncodeToString(expectedHash)

	cmd := &cli.ValidateCmd{
		File:     testFile,
		Expected: expectedHex,
		Algo:     "sha256",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestValidateCmd_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "validate.txt")
	testContent := "content to validate"

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	// Use wrong expected hash
	wrongHash := "0000000000000000000000000000000000000000000000000000000000000000"

	cmd := &cli.ValidateCmd{
		File:     testFile,
		Expected: wrongHash,
		Algo:     "sha256",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
}

func TestValidateCmd_NonexistentFile(t *testing.T) {
	cmd := &cli.ValidateCmd{
		File:     "/nonexistent/file.txt",
		Expected: "dummy-hash",
		Algo:     "sha256",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.Error(t, err)
}

func TestCompareCmd_EqualHashes(t *testing.T) {
	hash1 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hash2 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	cmd := &cli.CompareCmd{
		Hash1: hash1,
		Hash2: hash2,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestCompareCmd_DifferentHashes(t *testing.T) {
	hash1 := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	hash2 := "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"

	cmd := &cli.CompareCmd{
		Hash1: hash1,
		Hash2: hash2,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err) // Command succeeds but indicates hashes are different
}

func TestCompareCmd_InvalidHexString(t *testing.T) {
	tests := []struct {
		name  string
		hash1 string
		hash2 string
	}{
		{
			name:  "invalid first hash",
			hash1: "invalid-hex-string",
			hash2: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "invalid second hash",
			hash1: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			hash2: "not-valid-hex",
		},
		{
			name:  "both invalid",
			hash1: "invalid-hex1",
			hash2: "invalid-hex2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.CompareCmd{
				Hash1: tt.hash1,
				Hash2: tt.hash2,
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid")
		})
	}
}

func TestCompareCmd_WithWhitespace(t *testing.T) {
	// Test that whitespace is properly trimmed
	hash1 := "  2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824  "
	hash2 := "\t2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824\n"

	cmd := &cli.CompareCmd{
		Hash1: hash1,
		Hash2: hash2,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestCompareCmd_DifferentLengthHashes(t *testing.T) {
	// SHA256 vs MD5 (different lengths)
	sha256Hash := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	md5Hash := "5d41402abc4b2a76b9719d911017c592"

	cmd := &cli.CompareCmd{
		Hash1: sha256Hash,
		Hash2: md5Hash,
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err) // Should complete but indicate they're different
}

func TestHashFileCmd_PathSanitization(t *testing.T) {
	// Test that path sanitization works correctly
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "path sanitization test"

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	// Test with relative path that should be sanitized
	relPath := filepath.Base(testFile)
	
	// Change to the temp directory to test relative path resolution
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)
	
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	cmd := &cli.HashFileCmd{
		Path: relPath,
		Algo: "sha256",
	}
	ctx := newTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashDirCmd_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := &cli.HashDirCmd{
		Path: tmpDir,
		Algo: "sha256",
	}
	ctx := newTestContext()

	err := cmd.Run(ctx)
	require.NoError(t, err)
}

func TestHashStringCmd_UnicodeContent(t *testing.T) {
	unicodeTexts := []string{
		"🔐 Hash this emoji 🗝️",
		"Привет мир",
		"こんにちは世界", 
		"مرحبا بالعالم",
		"Mixed: Hello 世界 🌍",
	}

	for _, text := range unicodeTexts {
		t.Run("unicode_"+text[:10], func(t *testing.T) {
			cmd := &cli.HashStringCmd{
				Text: text,
				Algo: "sha256",
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestHMACCmd_EmptyInputs(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		key      string
		shouldErr bool
	}{
		{"empty text", "", "key", false},
		{"empty key", "text", "", false}, 
		{"both empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.HMACCmd{
				Text: tt.text,
				Key:  tt.key,
				Algo: "sha256",
			}
			ctx := newTestContext()

			err := cmd.Run(ctx)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}


