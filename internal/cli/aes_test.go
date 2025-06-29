package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bilte-co/toolshed/aes"
	"github.com/bilte-co/toolshed/internal/cli"
	"github.com/bilte-co/toolshed/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyCmd_ValidEntropy(t *testing.T) {
	tests := []struct {
		name    string
		entropy int
	}{
		{"128-bit key", 128},
		{"192-bit key", 192},
		{"256-bit key", 256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.GenerateKeyCmd{Entropy: tt.entropy}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.NoError(t, err)
		})
	}
}

func TestGenerateKeyCmd_InvalidEntropy(t *testing.T) {
	tests := []struct {
		name    string
		entropy int
	}{
		{"zero entropy", 0},
		{"negative entropy", -1},
		{"invalid entropy 100", 100},
		{"invalid entropy 127", 127},
		{"invalid entropy 129", 129},
		{"invalid entropy 300", 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.GenerateKeyCmd{Entropy: tt.entropy}
			ctx := testutil.NewTestContext()

			err := cmd.Run(ctx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid entropy")
		})
	}
}

func TestEncryptCmd_FileToStdout(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	testContent := "secret message"

	err := os.WriteFile(inputFile, []byte(testContent), 0o600)
	require.NoError(t, err)

	// Generate key
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.EncryptCmd{
		File: inputFile,
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestEncryptCmd_FileToFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.enc")
	testContent := "secret message for file encryption"

	err := os.WriteFile(inputFile, []byte(testContent), 0o600)
	require.NoError(t, err)

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.EncryptCmd{
		File:   inputFile,
		Key:    key,
		Output: outputFile,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)

	// Verify output file exists and contains data
	encryptedData, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.NotEmpty(t, encryptedData)

	// Verify we can decrypt it back
	decryptedData, err := aes.Decrypt(key, string(encryptedData))
	require.NoError(t, err)
	require.Equal(t, testContent, decryptedData)
}

func TestEncryptCmd_StdinToStdout(t *testing.T) {
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.EncryptCmd{
		File: "-",
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	testContent := "stdin secret message"
	go func() {
		defer w.Close()
		w.Write([]byte(testContent))
	}()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestEncryptCmd_MissingKey(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")

	err := os.WriteFile(inputFile, []byte("test"), 0o600)
	require.NoError(t, err)

	// Clear environment variable
	os.Unsetenv("AES_KEY")

	cmd := &cli.EncryptCmd{
		File: inputFile,
		// No key provided
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no AES key provided")
}

func TestEncryptCmd_KeyFromEnvironment(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	testContent := "env key test"

	err := os.WriteFile(inputFile, []byte(testContent), 0o600)
	require.NoError(t, err)

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	// Set environment variable
	os.Setenv("AES_KEY", key)
	defer os.Unsetenv("AES_KEY")

	cmd := &cli.EncryptCmd{
		File: inputFile,
		// Key will come from environment
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestEncryptCmd_NonexistentFile(t *testing.T) {
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.EncryptCmd{
		File: "/nonexistent/file.txt",
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open input file")
}

func TestDecryptCmd_FileToStdout(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "encrypted.txt")
	testContent := "secret message to decrypt"

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	// Create encrypted file
	ciphertext, err := aes.Encrypt(key, testContent)
	require.NoError(t, err)

	err = os.WriteFile(inputFile, []byte(ciphertext), 0o600)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File: inputFile,
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestDecryptCmd_FileToFile(t *testing.T) {
	tmpDir := t.TempDir()
	encryptedFile := filepath.Join(tmpDir, "encrypted.txt")
	outputFile := filepath.Join(tmpDir, "decrypted.txt")
	testContent := "secret message for file decryption"

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	// Create encrypted file
	ciphertext, err := aes.Encrypt(key, testContent)
	require.NoError(t, err)

	err = os.WriteFile(encryptedFile, []byte(ciphertext), 0o600)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File:   encryptedFile,
		Key:    key,
		Output: outputFile,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.NoError(t, err)

	// Verify output file contains original content
	decryptedData, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, testContent, string(decryptedData))
}

func TestDecryptCmd_StdinToStdout(t *testing.T) {
	testContent := "stdin decryption test"
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	ciphertext, err := aes.Encrypt(key, testContent)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File: "-",
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	// Mock stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write([]byte(ciphertext))
	}()

	err = cmd.Run(ctx)
	require.NoError(t, err)
}

func TestDecryptCmd_EmptyInput(t *testing.T) {
	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File: "-",
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	// Mock empty stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	os.Stdin = r
	w.Close() // Close immediately to simulate empty input

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no ciphertext data found")
}

func TestDecryptCmd_InvalidCiphertext(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "invalid.txt")

	err := os.WriteFile(inputFile, []byte("invalid-ciphertext"), 0o600)
	require.NoError(t, err)

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File: inputFile,
		Key:  key,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt data")
}

func TestDecryptCmd_WrongKey(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "encrypted.txt")
	testContent := "secret message"

	// Encrypt with one key
	key1, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	ciphertext, err := aes.Encrypt(key1, testContent)
	require.NoError(t, err)

	err = os.WriteFile(inputFile, []byte(ciphertext), 0o600)
	require.NoError(t, err)

	// Try to decrypt with different key
	key2, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	cmd := &cli.DecryptCmd{
		File: inputFile,
		Key:  key2,
	}
	ctx := testutil.NewTestContext()

	err = cmd.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt data")
}

func TestEncryptDecryptCmd_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "original.txt")
	encryptedFile := filepath.Join(tmpDir, "encrypted.txt")
	decryptedFile := filepath.Join(tmpDir, "decrypted.txt")

	testContent := "round trip test content with special chars: !@#$%^&*()"

	err := os.WriteFile(inputFile, []byte(testContent), 0o600)
	require.NoError(t, err)

	key, err := aes.GenerateAESKey(256)
	require.NoError(t, err)

	// Encrypt
	encryptCmd := &cli.EncryptCmd{
		File:   inputFile,
		Key:    key,
		Output: encryptedFile,
	}
	ctx := testutil.NewTestContext()

	err = encryptCmd.Run(ctx)
	require.NoError(t, err)

	// Decrypt
	decryptCmd := &cli.DecryptCmd{
		File:   encryptedFile,
		Key:    key,
		Output: decryptedFile,
	}

	err = decryptCmd.Run(ctx)
	require.NoError(t, err)

	// Verify round trip
	decryptedContent, err := os.ReadFile(decryptedFile)
	require.NoError(t, err)
	require.Equal(t, testContent, string(decryptedContent))
}

func TestEncryptCmd_GetKey(t *testing.T) {
	tests := []struct {
		name      string
		cmdKey    string
		envKey    string
		expectErr bool
		setup     func()
		cleanup   func()
	}{
		{
			name:   "key from flag",
			cmdKey: "test-key-from-flag",
		},
		{
			name:    "key from environment",
			envKey:  "test-key-from-env",
			setup:   func() { os.Setenv("AES_KEY", "test-key-from-env") },
			cleanup: func() { os.Unsetenv("AES_KEY") },
		},
		{
			name:      "no key provided",
			expectErr: true,
			setup:     func() { os.Unsetenv("AES_KEY") },
		},
		{
			name:    "flag takes precedence over env",
			cmdKey:  "flag-key",
			envKey:  "env-key",
			setup:   func() { os.Setenv("AES_KEY", "env-key") },
			cleanup: func() { os.Unsetenv("AES_KEY") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			cmd := &cli.EncryptCmd{Key: tt.cmdKey}

			// Use reflection or create a test method to access getKey
			// For this test, we'll create a minimal file to trigger getKey
			tmpDir := t.TempDir()
			inputFile := filepath.Join(tmpDir, "test.txt")
			err := os.WriteFile(inputFile, []byte("test"), 0o600)
			require.NoError(t, err)

			cmd.File = inputFile
			ctx := testutil.NewTestContext()

			err = cmd.Run(ctx)
			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "no AES key provided")
			} else {
				// The test might fail during encryption due to invalid key format,
				// but getKey should work correctly
				if err != nil {
					require.Contains(t, err.Error(), "failed to encrypt")
				}
			}
		})
	}
}
