package hash

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestZipFile(t *testing.T, files map[string]string) string {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")

	file, err := os.Create(zipPath)
	require.NoError(t, err)
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	for filename, content := range files {
		writer, err := zipWriter.Create(filename)
		require.NoError(t, err)
		_, err = writer.Write([]byte(content))
		require.NoError(t, err)
	}

	return zipPath
}

func createTestTarFile(t *testing.T, files map[string]string) string {
	tmpDir := t.TempDir()
	tarPath := filepath.Join(tmpDir, "test.tar")

	file, err := os.Create(tarPath)
	require.NoError(t, err)
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	for filename, content := range files {
		header := &tar.Header{
			Name: filename,
			Size: int64(len(content)),
			Mode: 0644,
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(t, err)
		_, err = tarWriter.Write([]byte(content))
		require.NoError(t, err)
	}

	return tarPath
}

func createTestTarGzFile(t *testing.T, files map[string]string) string {
	tmpDir := t.TempDir()
	targzPath := filepath.Join(tmpDir, "test.tar.gz")

	file, err := os.Create(targzPath)
	require.NoError(t, err)
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	for filename, content := range files {
		header := &tar.Header{
			Name: filename,
			Size: int64(len(content)),
			Mode: 0644,
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(t, err)
		_, err = tarWriter.Write([]byte(content))
		require.NoError(t, err)
	}

	return targzPath
}

func createTestGzFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	gzPath := filepath.Join(tmpDir, "test.txt.gz")

	file, err := os.Create(gzPath)
	require.NoError(t, err)
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	_, err = gzWriter.Write([]byte(content))
	require.NoError(t, err)

	return gzPath
}

func TestHashArchive_ZIP(t *testing.T) {
	testFiles := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	zipPath := createTestZipFile(t, testFiles)

	// Test ZIP hashing
	hash, err := HashArchive(zipPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test deterministic behavior
	hash2, err := HashArchive(zipPath, "sha256")
	require.NoError(t, err)
	assert.Equal(t, hash, hash2, "ZIP hashing should be deterministic")

	// Test different algorithm
	hash3, err := HashArchive(zipPath, "sha512")
	require.NoError(t, err)
	assert.NotEqual(t, hash, hash3, "Different algorithms should produce different hashes")
}

func TestHashArchive_TAR(t *testing.T) {
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	tarPath := createTestTarFile(t, testFiles)

	// Test TAR hashing
	hash, err := HashArchive(tarPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test deterministic behavior
	hash2, err := HashArchive(tarPath, "sha256")
	require.NoError(t, err)
	assert.Equal(t, hash, hash2, "TAR hashing should be deterministic")
}

func TestHashArchive_TARGZ(t *testing.T) {
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	targzPath := createTestTarGzFile(t, testFiles)

	// Test TAR.GZ hashing
	hash, err := HashArchive(targzPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test deterministic behavior
	hash2, err := HashArchive(targzPath, "sha256")
	require.NoError(t, err)
	assert.Equal(t, hash, hash2, "TAR.GZ hashing should be deterministic")
}

func TestHashArchive_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	unsupportedFile := filepath.Join(tmpDir, "test.rar")
	err := os.WriteFile(unsupportedFile, []byte("dummy content"), 0644)
	require.NoError(t, err)

	_, err = HashArchive(unsupportedFile, "sha256")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported archive format")
}

func TestHashArchive_NonExistentFile(t *testing.T) {
	_, err := HashArchive("/nonexistent/archive.zip", "sha256")
	assert.Error(t, err)
}

func TestHashArchive_InvalidAlgorithm(t *testing.T) {
	testFiles := map[string]string{
		"file1.txt": "content1",
	}
	zipPath := createTestZipFile(t, testFiles)

	_, err := HashArchive(zipPath, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported hash algorithm")
}

func TestHashArchive_EmptyArchive(t *testing.T) {
	// Test with empty ZIP
	emptyFiles := map[string]string{}
	zipPath := createTestZipFile(t, emptyFiles)

	hash, err := HashArchive(zipPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash) // Should still produce a hash even for empty archive
}

func TestHashArchiveWithOptions(t *testing.T) {
	testFiles := map[string]string{
		"file1.txt": "content1",
	}
	zipPath := createTestZipFile(t, testFiles)

	// Test hex format
	result, err := HashArchiveWithOptions(zipPath, "sha256", Options{Format: FormatHex})
	require.NoError(t, err)
	assert.IsType(t, "", result)

	// Test base64 format
	result, err = HashArchiveWithOptions(zipPath, "sha256", Options{Format: FormatBase64})
	require.NoError(t, err)
	assert.IsType(t, "", result)

	// Test raw format
	result, err = HashArchiveWithOptions(zipPath, "sha256", Options{Format: FormatRaw})
	require.NoError(t, err)
	assert.IsType(t, []byte{}, result)

	// Test with prefix
	result, err = HashArchiveWithOptions(zipPath, "sha256", Options{Format: FormatHex, Prefix: true})
	require.NoError(t, err)
	resultStr := result.(string)
	assert.Contains(t, resultStr, "sha256:")
}

func TestHashCompressedFile(t *testing.T) {
	testContent := "hello world compressed"
	gzPath := createTestGzFile(t, testContent)

	// Test compressed file hashing
	hash, err := HashCompressedFile(gzPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Compare with direct hashing of uncompressed content
	expectedHash, err := HashString(testContent, "sha256")
	require.NoError(t, err)
	assert.Equal(t, expectedHash, hash, "Compressed file hash should match uncompressed content")
}

func TestHashCompressedFile_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	unsupportedFile := filepath.Join(tmpDir, "test.bz2")
	err := os.WriteFile(unsupportedFile, []byte("dummy content"), 0644)
	require.NoError(t, err)

	_, err = HashCompressedFile(unsupportedFile, "sha256")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported compressed file format")
}

func TestHashCompressedFile_NonExistentFile(t *testing.T) {
	_, err := HashCompressedFile("/nonexistent/file.gz", "sha256")
	assert.Error(t, err)
}

func TestHashCompressedFile_InvalidGzip(t *testing.T) {
	tmpDir := t.TempDir()
	invalidGzFile := filepath.Join(tmpDir, "invalid.gz")
	err := os.WriteFile(invalidGzFile, []byte("not a gzip file"), 0644)
	require.NoError(t, err)

	_, err = HashCompressedFile(invalidGzFile, "sha256")
	assert.Error(t, err)
}

func TestHashCompressedFileWithOptions(t *testing.T) {
	testContent := "test content"
	gzPath := createTestGzFile(t, testContent)

	// Test hex format
	result, err := HashCompressedFileWithOptions(gzPath, "sha256", Options{Format: FormatHex})
	require.NoError(t, err)
	assert.IsType(t, "", result)

	// Test with prefix
	result, err = HashCompressedFileWithOptions(gzPath, "sha256", Options{Format: FormatHex, Prefix: true})
	require.NoError(t, err)
	resultStr := result.(string)
	assert.Contains(t, resultStr, "sha256:")
}

func TestHashArchive_DeterministicOrdering(t *testing.T) {
	// Test that files are processed in deterministic order regardless of creation order
	testFiles := map[string]string{
		"z_file.txt": "content_z",
		"a_file.txt": "content_a",
		"m_file.txt": "content_m",
	}

	// Create two identical archives
	zipPath1 := createTestZipFile(t, testFiles)
	zipPath2 := createTestZipFile(t, testFiles)

	hash1, err := HashArchive(zipPath1, "sha256")
	require.NoError(t, err)

	hash2, err := HashArchive(zipPath2, "sha256")
	require.NoError(t, err)

	assert.Equal(t, hash1, hash2, "Archives with same content should have same hash regardless of creation order")
}

func TestHashArchive_WithDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")

	file, err := os.Create(zipPath)
	require.NoError(t, err)
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Add a directory entry (should be skipped)
	_, err = zipWriter.Create("testdir/")
	require.NoError(t, err)

	// Add a file
	writer, err := zipWriter.Create("testdir/file.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content"))
	require.NoError(t, err)

	zipWriter.Close()
	file.Close()

	// Should successfully hash, skipping directory entries
	hash, err := HashArchive(zipPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashArchive_LargeFiles(t *testing.T) {
	// Test with moderately large content to ensure proper handling
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	testFiles := map[string]string{
		"large_file.bin": string(largeContent),
	}

	zipPath := createTestZipFile(t, testFiles)

	hash, err := HashArchive(zipPath, "sha256")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashArchive_AllAlgorithms(t *testing.T) {
	testFiles := map[string]string{
		"test.txt": "test content",
	}
	zipPath := createTestZipFile(t, testFiles)

	algorithms := []string{"md5", "sha1", "sha256", "sha512", "blake2b"}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			hash, err := HashArchive(zipPath, algo)
			require.NoError(t, err, "Algorithm %s should work", algo)
			assert.NotEmpty(t, hash)
		})
	}
}

func BenchmarkHashArchive_ZIP(b *testing.B) {
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3",
	}

	tmpDir := b.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")

	file, err := os.Create(zipPath)
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	for filename, content := range testFiles {
		writer, err := zipWriter.Create(filename)
		if err != nil {
			b.Fatal(err)
		}
		_, err = writer.Write([]byte(content))
		if err != nil {
			b.Fatal(err)
		}
	}
	zipWriter.Close()
	file.Close()

	for b.Loop() {
		_, err := HashArchive(zipPath, "sha256")
		if err != nil {
			b.Fatal(err)
		}
	}
}
