package hash

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HashArchive hashes the contents of an archive file (.zip, .tar.gz, .tar).
// The hash is computed over the sorted list of file entries and their contents
// to ensure deterministic results.
func HashArchive(path string, algorithm string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch {
	case ext == ".zip":
		return hashZipArchive(path, algorithm)
	case ext == ".gz" && strings.HasSuffix(strings.ToLower(path), ".tar.gz"):
		return hashTarGzArchive(path, algorithm)
	case ext == ".tar":
		return hashTarArchive(path, algorithm)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// HashArchiveWithOptions hashes an archive with custom options.
func HashArchiveWithOptions(path string, algorithm string, opts Options) (interface{}, error) {
	data, err := HashArchive(path, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}

// archiveEntry represents a file entry in an archive.
type archiveEntry struct {
	name string
	data []byte
}

// hashZipArchive hashes a ZIP archive.
func hashZipArchive(path string, algorithm string) ([]byte, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP archive %s: %w", path, err)
	}
	defer reader.Close()
	
	var entries []archiveEntry
	
	for _, file := range reader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}
		
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s in ZIP archive: %w", file.Name, err)
		}
		
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s in ZIP archive: %w", file.Name, err)
		}
		
		entries = append(entries, archiveEntry{
			name: file.Name,
			data: data,
		})
	}
	
	return hashArchiveEntries(entries, algorithm)
}

// hashTarGzArchive hashes a compressed TAR archive.
func hashTarGzArchive(path string, algorithm string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar.gz archive %s: %w", path, err)
	}
	defer file.Close()
	
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader for %s: %w", path, err)
	}
	defer gzReader.Close()
	
	return hashTarReader(gzReader, algorithm)
}

// hashTarArchive hashes a TAR archive.
func hashTarArchive(path string, algorithm string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar archive %s: %w", path, err)
	}
	defer file.Close()
	
	return hashTarReader(file, algorithm)
}

// hashTarReader hashes a TAR archive from an io.Reader.
func hashTarReader(reader io.Reader, algorithm string) ([]byte, error) {
	tarReader := tar.NewReader(reader)
	var entries []archiveEntry
	
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}
		
		// Skip directories and other non-regular files
		if header.Typeflag != tar.TypeReg {
			continue
		}
		
		data, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s in tar archive: %w", header.Name, err)
		}
		
		entries = append(entries, archiveEntry{
			name: header.Name,
			data: data,
		})
	}
	
	return hashArchiveEntries(entries, algorithm)
}

// hashArchiveEntries hashes a sorted list of archive entries.
func hashArchiveEntries(entries []archiveEntry, algorithm string) ([]byte, error) {
	h, err := getHasher(algorithm)
	if err != nil {
		return nil, err
	}
	
	// Sort entries by name for deterministic output
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})
	
	// Hash each entry's name and content
	for _, entry := range entries {
		h.Write([]byte(entry.name))
		h.Write(entry.data)
	}
	
	return h.Sum(nil), nil
}

// HashCompressedFile hashes a compressed file by first decompressing it.
// Supports .gz files.
func HashCompressedFile(path string, algorithm string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".gz":
		return hashGzipFile(path, algorithm)
	default:
		return nil, fmt.Errorf("unsupported compressed file format: %s", ext)
	}
}

// hashGzipFile hashes a gzip-compressed file.
func hashGzipFile(path string, algorithm string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip file %s: %w", path, err)
	}
	defer file.Close()
	
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader for %s: %w", path, err)
	}
	defer gzReader.Close()
	
	return HashReader(gzReader, algorithm)
}

// HashCompressedFileWithOptions hashes a compressed file with custom options.
func HashCompressedFileWithOptions(path string, algorithm string, opts Options) (interface{}, error) {
	data, err := HashCompressedFile(path, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}
