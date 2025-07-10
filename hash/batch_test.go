package hash

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashFilesInParallel(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3",
	}

	var filePaths []string
	for filename, content := range testFiles {
		path := filepath.Join(tmpDir, filename)
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err, "Failed to create test file")
		filePaths = append(filePaths, path)
	}

	// Test parallel hashing
	result := HashFilesInParallel(filePaths, "sha256", 2)

	// Should have results for all files
	assert.Len(t, result.Results, len(filePaths), "Should have result for each file")
	assert.Empty(t, result.Errors, "Should have no errors")

	// All results should have hashes
	for _, fileResult := range result.Results {
		assert.NotEmpty(t, fileResult.Hash, "Hash should not be empty")
		assert.Equal(t, "sha256", fileResult.Algorithm, "Algorithm should match")
		assert.NoError(t, fileResult.Error, "File result should have no error")
		assert.Contains(t, filePaths, fileResult.Path, "Path should be in original list")
	}

	// Test with different worker counts
	result2 := HashFilesInParallel(filePaths, "sha256", 1)
	assert.Len(t, result2.Results, len(filePaths), "Should work with 1 worker")

	result3 := HashFilesInParallel(filePaths, "sha256", 0) // Should default to CPU count
	assert.Len(t, result3.Results, len(filePaths), "Should work with 0 workers (default)")
}

func TestHashFilesInParallelWithErrors(t *testing.T) {
	// Include non-existent file
	filePaths := []string{
		"/nonexistent/file1.txt",
		"/nonexistent/file2.txt",
	}

	result := HashFilesInParallel(filePaths, "sha256", 2)

	// Should have results for all files (with errors)
	assert.Len(t, result.Results, len(filePaths), "Should have result for each file")
	assert.Len(t, result.Errors, len(filePaths), "Should have errors for all files")

	// All results should have errors
	for _, fileResult := range result.Results {
		assert.Error(t, fileResult.Error, "File result should have error for non-existent file")
		assert.Empty(t, fileResult.Hash, "Hash should be empty on error")
	}
}

func TestHashFilesInParallelWithOptions(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err, "Failed to create test file")

	filePaths := []string{testFile}

	// Test with hex format
	result := HashFilesInParallelWithOptions(filePaths, "sha256", 1, Options{Format: FormatHex})
	require.Len(t, result.Results, 1, "Should have one result")
	require.Empty(t, result.Errors, "Should have no errors")

	// Hash should be formatted as hex string (converted to bytes in our implementation)
	assert.NotEmpty(t, result.Results[0].Hash, "Hash should not be empty")
}

func TestValidateFileChecksum(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "hello world"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Failed to create test file")

	// Get the correct hash
	expectedHash, err := HashFile(testFile, "sha256")
	require.NoError(t, err, "Failed to hash test file")
	expectedHexHash := hex.EncodeToString(expectedHash)

	// Test successful validation
	err = ValidateFileChecksum(testFile, expectedHexHash, "sha256")
	assert.NoError(t, err, "Validation should succeed with correct hash")

	// Test with algorithm prefix
	err = ValidateFileChecksum(testFile, "sha256:"+expectedHexHash, "sha256")
	assert.NoError(t, err, "Validation should succeed with algorithm prefix")

	// Test with wrong hash
	err = ValidateFileChecksum(testFile, "0123456789abcdef", "sha256")
	assert.Error(t, err, "Validation should fail with wrong hash")
	assert.Contains(t, err.Error(), "checksum mismatch", "Error should mention checksum mismatch")

	// Test with invalid hash format
	err = ValidateFileChecksum(testFile, "invalid-hex", "sha256")
	assert.Error(t, err, "Validation should fail with invalid hash format")

	// Test with non-existent file
	err = ValidateFileChecksum("/nonexistent/file.txt", expectedHexHash, "sha256")
	assert.Error(t, err, "Validation should fail for non-existent file")
}

func TestValidateFilesInParallel(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3",
	}

	var checksums []FileChecksum
	for filename, content := range testFiles {
		path := filepath.Join(tmpDir, filename)
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err, "Failed to create test file")

		// Get correct hash
		hash, err := HashFile(path, "sha256")
		require.NoError(t, err, "Failed to hash test file")

		checksums = append(checksums, FileChecksum{
			Path:         path,
			ExpectedHash: hex.EncodeToString(hash),
		})
	}

	// Test successful validation
	errors := ValidateFilesInParallel(checksums, "sha256", 2)
	assert.Empty(t, errors, "Should have no validation errors")

	// Test with one wrong checksum
	checksums[0].ExpectedHash = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	errors = ValidateFilesInParallel(checksums, "sha256", 2)
	assert.Len(t, errors, 1, "Should have one validation error")
	assert.Contains(t, errors[0].Error(), "checksum mismatch", "Error should mention checksum mismatch")
}

func TestEmptyBatchOperations(t *testing.T) {
	// Test with empty file list
	result := HashFilesInParallel([]string{}, "sha256", 2)
	assert.Empty(t, result.Results, "Should have no results for empty list")
	assert.Empty(t, result.Errors, "Should have no errors for empty list")

	// Test validation with empty list
	errors := ValidateFilesInParallel([]FileChecksum{}, "sha256", 2)
	assert.Empty(t, errors, "Should have no errors for empty validation list")
}

func TestHashFilesInParallel_UnsupportedAlgorithm(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	result := HashFilesInParallel([]string{testFile}, "unsupported", 1)
	assert.Len(t, result.Results, 1)
	assert.Len(t, result.Errors, 1)
	assert.Error(t, result.Results[0].Error)
	assert.Contains(t, result.Results[0].Error.Error(), "unsupported hash algorithm")
}

func TestHashFilesInParallel_NegativeWorkers(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	// Should default to CPU count
	result := HashFilesInParallel([]string{testFile}, "sha256", -5)
	assert.Len(t, result.Results, 1)
	assert.Empty(t, result.Errors)
	assert.NoError(t, result.Results[0].Error)
}

func TestHashFilesInParallelWithOptions_UnsupportedAlgorithm(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	result := HashFilesInParallelWithOptions([]string{testFile}, "unsupported", 1, Options{Format: FormatHex})
	assert.Len(t, result.Results, 1)
	assert.Len(t, result.Errors, 1)
	assert.Error(t, result.Results[0].Error)
}

func TestHashFilesInParallelWithOptions_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	result := HashFilesInParallelWithOptions([]string{testFile}, "sha256", 1, Options{Format: Format("invalid")})
	assert.Len(t, result.Results, 1)
	assert.Len(t, result.Errors, 1)
	assert.Error(t, result.Results[0].Error)
	assert.Contains(t, result.Errors[0].Error(), "failed to format output")
}

func TestHashFilesInParallelWithOptions_FormatTypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	// Test different formats
	formats := []Format{FormatHex, FormatBase64, FormatRaw}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			result := HashFilesInParallelWithOptions([]string{testFile}, "sha256", 1, Options{Format: format})
			assert.Len(t, result.Results, 1)
			assert.Empty(t, result.Errors)
			assert.NoError(t, result.Results[0].Error)
			assert.NotEmpty(t, result.Results[0].Hash)
		})
	}
}

func TestValidateFileChecksum_EmptyHash(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	err = ValidateFileChecksum(testFile, "", "sha256")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")
}

func TestValidateFileChecksum_ShortHash(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	err = ValidateFileChecksum(testFile, "abc", "sha256")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid hash format")
}

func TestValidateFileChecksum_UnsupportedAlgorithm(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	err = ValidateFileChecksum(testFile, "abc123", "unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash file")
}

func TestValidateFilesInParallel_EdgeCases(t *testing.T) {
	// Test with negative workers
	errors := ValidateFilesInParallel([]FileChecksum{}, "sha256", -1)
	assert.Empty(t, errors)

	// Test with zero workers
	errors = ValidateFilesInParallel([]FileChecksum{}, "sha256", 0)
	assert.Empty(t, errors)
}

func TestValidateFilesInParallel_MixedResults(t *testing.T) {
	tmpDir := t.TempDir()

	// Create one valid file
	validFile := filepath.Join(tmpDir, "valid.txt")
	err := os.WriteFile(validFile, []byte("valid content"), 0644)
	require.NoError(t, err)

	validHash, err := HashFile(validFile, "sha256")
	require.NoError(t, err)

	checksums := []FileChecksum{
		{
			Path:         validFile,
			ExpectedHash: hex.EncodeToString(validHash),
		},
		{
			Path:         "/nonexistent/file.txt",
			ExpectedHash: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		},
		{
			Path:         validFile,
			ExpectedHash: "wronghash123", // Invalid hash format
		},
	}

	errors := ValidateFilesInParallel(checksums, "sha256", 2)
	assert.Len(t, errors, 2, "Should have 2 errors: nonexistent file and invalid hash")
}

func TestFileHashResult_Fields(t *testing.T) {
	// Test that FileHashResult has all expected fields
	result := FileHashResult{
		Path:      "/test/path",
		Hash:      []byte{0x01, 0x02, 0x03},
		Error:     nil,
		Algorithm: "sha256",
	}

	assert.Equal(t, "/test/path", result.Path)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result.Hash)
	assert.NoError(t, result.Error)
	assert.Equal(t, "sha256", result.Algorithm)
}

func TestBatchHashResult_Fields(t *testing.T) {
	// Test that BatchHashResult has all expected fields
	result := BatchHashResult{
		Results: []FileHashResult{
			{Path: "test1", Hash: []byte{0x01}, Algorithm: "sha256"},
			{Path: "test2", Hash: []byte{0x02}, Algorithm: "sha256"},
		},
		Errors: []error{
			fmt.Errorf("test error"),
		},
	}

	assert.Len(t, result.Results, 2)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "test1", result.Results[0].Path)
	assert.Equal(t, "test2", result.Results[1].Path)
	assert.Contains(t, result.Errors[0].Error(), "test error")
}

func TestFileChecksum_Fields(t *testing.T) {
	// Test that FileChecksum has all expected fields
	checksum := FileChecksum{
		Path:         "/test/path",
		ExpectedHash: "abc123",
	}

	assert.Equal(t, "/test/path", checksum.Path)
	assert.Equal(t, "abc123", checksum.ExpectedHash)
}

func TestHashFilesInParallel_LargeNumberOfFiles(t *testing.T) {
	tmpDir := t.TempDir()

	var filePaths []string
	numFiles := 100

	// Create many small files
	for i := range numFiles {
		path := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("content%d", i)
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)
		filePaths = append(filePaths, path)
	}

	// Test with different worker counts
	workerCounts := []int{1, 4, 10, 20}

	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers_%d", workers), func(t *testing.T) {
			result := HashFilesInParallel(filePaths, "sha256", workers)
			assert.Len(t, result.Results, numFiles)
			assert.Empty(t, result.Errors)

			// All results should have hashes
			for _, fileResult := range result.Results {
				assert.NoError(t, fileResult.Error)
				assert.NotEmpty(t, fileResult.Hash)
				assert.Equal(t, "sha256", fileResult.Algorithm)
			}
		})
	}
}

func TestValidateFileChecksum_AlgorithmPrefixEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "test content"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	expectedHash, err := HashFile(testFile, "sha256")
	require.NoError(t, err)
	expectedHexHash := hex.EncodeToString(expectedHash)

	// Test with prefix but different algorithm
	err = ValidateFileChecksum(testFile, "sha512:"+expectedHexHash, "sha256")
	assert.Error(t, err) // Should fail because we're stripping sha512: but using sha256

	// Test with malformed prefix
	err = ValidateFileChecksum(testFile, "sha256"+expectedHexHash, "sha256") // Missing colon
	assert.Error(t, err)                                                     // Should fail because hash format is wrong
}

func BenchmarkHashFilesInParallel(b *testing.B) {
	// Create temporary test files
	tmpDir := b.TempDir()

	var filePaths []string
	for i := range 10 {
		path := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("content for file %d", i)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			b.Fatal(err)
		}
		filePaths = append(filePaths, path)
	}

	for b.Loop() {
		result := HashFilesInParallel(filePaths, "sha256", 4)
		if len(result.Errors) > 0 {
			b.Fatal("Unexpected errors in benchmark")
		}
	}
}
