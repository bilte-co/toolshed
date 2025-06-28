// Package hash provides a flexible and extensible API for hashing operations.
// It supports multiple algorithms including SHA-1, SHA-256, SHA-512, BLAKE2b, and MD5
// with security features like HMAC, password hashing, and constant-time comparison.
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/crypto/blake2b"
)

var (
	// ErrUnsupportedAlgorithm is returned when an unsupported hash algorithm is requested.
	ErrUnsupportedAlgorithm = errors.New("unsupported hash algorithm")

	// ErrInvalidFormat is returned when an invalid output format is requested.
	ErrInvalidFormat = errors.New("invalid output format")

	// customHashers stores registered custom hash algorithms.
	customHashers = make(map[string]func() hash.Hash)
	hasherMutex   sync.RWMutex
)

// Format represents the output format for hash results.
type Format string

const (
	// FormatRaw returns raw bytes.
	FormatRaw Format = "raw"
	// FormatHex returns hexadecimal encoding.
	FormatHex Format = "hex"
	// FormatBase64 returns base64 encoding.
	FormatBase64 Format = "base64"
)

// Options configures hash operations.
type Options struct {
	Format     Format
	Prefix     bool
	Workers    int
	BufferSize int
}

// DefaultOptions provides sensible defaults for hash operations.
var DefaultOptions = Options{
	Format:     FormatHex,
	Prefix:     false,
	Workers:    4,
	BufferSize: 64 * 1024, // 64KB
}

// Hasher wraps a hash.Hash with additional functionality.
type Hasher struct {
	hash.Hash
	algorithm string
}

// NewHasher creates a new Hasher instance for the specified algorithm.
func NewHasher(algorithm string) (*Hasher, error) {
	h, err := getHasher(algorithm)
	if err != nil {
		return nil, err
	}

	return &Hasher{
		Hash:      h,
		algorithm: algorithm,
	}, nil
}

// Algorithm returns the hash algorithm name.
func (h *Hasher) Algorithm() string {
	return h.algorithm
}

// SumHex returns the hex-encoded hash sum.
func (h *Hasher) SumHex() string {
	return hex.EncodeToString(h.Sum(nil))
}

// SumBase64 returns the base64-encoded hash sum.
func (h *Hasher) SumBase64() string {
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// getHasher returns a hash.Hash instance for the specified algorithm.
func getHasher(algorithm string) (hash.Hash, error) {
	algorithm = strings.ToLower(algorithm)

	// Warn about insecure algorithms
	if algorithm == "md5" || algorithm == "sha1" {
		log.Printf("WARNING: Using insecure hash algorithm %s. Consider using SHA-256 or SHA-512 instead.", algorithm)
	}

	switch algorithm {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	case "blake2b":
		h, err := blake2b.New256(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create BLAKE2b hasher: %w", err)
		}
		return h, nil
	default:
		// Check custom hashers
		hasherMutex.RLock()
		factory, exists := customHashers[algorithm]
		hasherMutex.RUnlock()

		if exists {
			return factory(), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}
}

// RegisterHasher registers a custom hash algorithm.
func RegisterHasher(name string, factory func() hash.Hash) {
	hasherMutex.Lock()
	defer hasherMutex.Unlock()
	customHashers[strings.ToLower(name)] = factory
}

// formatOutput formats the hash bytes according to the specified format and options.
func formatOutput(data []byte, algorithm string, opts Options) (any, error) {
	var result string

	switch opts.Format {
	case FormatRaw:
		return data, nil
	case FormatHex:
		result = hex.EncodeToString(data)
	case FormatBase64:
		result = base64.StdEncoding.EncodeToString(data)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, opts.Format)
	}

	if opts.Prefix {
		result = algorithm + ":" + result
	}

	return result, nil
}

// HashString hashes a string using the specified algorithm.
func HashString(input string, algorithm string) ([]byte, error) {
	return HashBytes([]byte(input), algorithm)
}

// HashStringWithOptions hashes a string with custom options.
func HashStringWithOptions(input string, algorithm string, opts Options) (any, error) {
	data, err := HashBytes([]byte(input), algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}

// HashBytes hashes a byte slice using the specified algorithm.
func HashBytes(input []byte, algorithm string) ([]byte, error) {
	h, err := getHasher(algorithm)
	if err != nil {
		return nil, err
	}

	h.Write(input)
	return h.Sum(nil), nil
}

// HashBytesWithOptions hashes bytes with custom options.
func HashBytesWithOptions(input []byte, algorithm string, opts Options) (any, error) {
	data, err := HashBytes(input, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}

// HashReader hashes data from an io.Reader using the specified algorithm.
func HashReader(r io.Reader, algorithm string) ([]byte, error) {
	h, err := getHasher(algorithm)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return h.Sum(nil), nil
}

// HashReaderWithOptions hashes an io.Reader with custom options.
func HashReaderWithOptions(r io.Reader, algorithm string, opts Options) (any, error) {
	data, err := HashReader(r, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}

// HashFile hashes a file using the specified algorithm.
func HashFile(path string, algorithm string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	return HashReader(file, algorithm)
}

// HashFileWithOptions hashes a file with custom options.
func HashFileWithOptions(path string, algorithm string, opts Options) (any, error) {
	data, err := HashFile(path, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}

// HashDir hashes a directory's contents deterministically.
func HashDir(path string, algorithm string, recursive bool) ([]byte, error) {
	h, err := getHasher(algorithm)
	if err != nil {
		return nil, err
	}

	var files []string

	walkFn := func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !recursive && filePath != path {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, filePath)
		return nil
	}

	if err := filepath.WalkDir(path, walkFn); err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", path, err)
	}

	// Sort files for deterministic output
	sort.Strings(files)

	for _, file := range files {
		relPath, err := filepath.Rel(path, file)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for %s: %w", file, err)
		}

		// Write file path to hash
		h.Write([]byte(relPath))

		// Hash file content
		fileData, err := HashFile(file, algorithm)
		if err != nil {
			return nil, fmt.Errorf("failed to hash file %s: %w", file, err)
		}
		h.Write(fileData)
	}

	return h.Sum(nil), nil
}

// HashDirWithOptions hashes a directory with custom options.
func HashDirWithOptions(path string, algorithm string, recursive bool, opts Options) (any, error) {
	data, err := HashDir(path, algorithm, recursive)
	if err != nil {
		return nil, err
	}
	return formatOutput(data, algorithm, opts)
}
