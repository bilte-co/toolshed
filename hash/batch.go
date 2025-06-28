package hash

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
)

// FileHashResult represents the result of hashing a single file.
type FileHashResult struct {
	Path      string
	Hash      []byte
	Error     error
	Algorithm string
}

// BatchHashResult represents the results of batch hashing operations.
type BatchHashResult struct {
	Results []FileHashResult
	Errors  []error
}

// HashFilesInParallel hashes multiple files in parallel using the specified number of workers.
// If workers is 0 or negative, it defaults to the number of CPU cores.
func HashFilesInParallel(paths []string, algorithm string, workers int) *BatchHashResult {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	
	if len(paths) == 0 {
		return &BatchHashResult{Results: []FileHashResult{}, Errors: []error{}}
	}
	
	// Create channels for work distribution and result collection
	jobs := make(chan string, len(paths))
	results := make(chan FileHashResult, len(paths))
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				hash, err := HashFile(path, algorithm)
				results <- FileHashResult{
					Path:      path,
					Hash:      hash,
					Error:     err,
					Algorithm: algorithm,
				}
			}
		}()
	}
	
	// Send jobs to workers
	go func() {
		defer close(jobs)
		for _, path := range paths {
			jobs <- path
		}
	}()
	
	// Wait for all workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results
	var fileResults []FileHashResult
	var errors []error
	
	for result := range results {
		fileResults = append(fileResults, result)
		if result.Error != nil {
			errors = append(errors, fmt.Errorf("failed to hash %s: %w", result.Path, result.Error))
		}
	}
	
	return &BatchHashResult{
		Results: fileResults,
		Errors:  errors,
	}
}

// HashFilesInParallelWithOptions hashes multiple files in parallel with custom options.
func HashFilesInParallelWithOptions(paths []string, algorithm string, workers int, opts Options) *BatchHashResult {
	batchResult := HashFilesInParallel(paths, algorithm, workers)
	
	// Format the results according to options
	for i := range batchResult.Results {
		if batchResult.Results[i].Error == nil {
			formatted, err := formatOutput(batchResult.Results[i].Hash, algorithm, opts)
			if err != nil {
				batchResult.Results[i].Error = err
				batchResult.Errors = append(batchResult.Errors, 
					fmt.Errorf("failed to format output for %s: %w", batchResult.Results[i].Path, err))
			} else {
				// Store formatted result in Hash field as interface{}
				switch f := formatted.(type) {
				case []byte:
					batchResult.Results[i].Hash = f
				case string:
					batchResult.Results[i].Hash = []byte(f)
				}
			}
		}
	}
	
	return batchResult
}

// ValidateFileChecksum validates a file against a known hash string.
func ValidateFileChecksum(path string, expectedHash string, algorithm string) error {
	actualHash, err := HashFile(path, algorithm)
	if err != nil {
		return fmt.Errorf("failed to hash file %s: %w", path, err)
	}
	
	// Parse expected hash (remove algorithm prefix if present)
	expected := expectedHash
	if prefix := algorithm + ":"; len(expectedHash) > len(prefix) && expectedHash[:len(prefix)] == prefix {
		expected = expectedHash[len(prefix):]
	}
	
	// Convert expected hash from hex to bytes
	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return fmt.Errorf("invalid hash format: %w", err)
	}
	
	// Compare hashes using constant-time comparison
	if !EqualConstantTime(actualHash, expectedBytes) {
		return fmt.Errorf("checksum mismatch for file %s: expected %s, got %s", 
			path, expected, hex.EncodeToString(actualHash))
	}
	
	return nil
}

// ValidateFilesChecksum validates multiple files against their expected checksums.
type FileChecksum struct {
	Path         string
	ExpectedHash string
}

// ValidateFilesInParallel validates multiple files against their checksums in parallel.
func ValidateFilesInParallel(checksums []FileChecksum, algorithm string, workers int) []error {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	
	if len(checksums) == 0 {
		return []error{}
	}
	
	// Create channels for work distribution and result collection
	jobs := make(chan FileChecksum, len(checksums))
	results := make(chan error, len(checksums))
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for checksum := range jobs {
				err := ValidateFileChecksum(checksum.Path, checksum.ExpectedHash, algorithm)
				results <- err
			}
		}()
	}
	
	// Send jobs to workers
	go func() {
		defer close(jobs)
		for _, checksum := range checksums {
			jobs <- checksum
		}
	}()
	
	// Wait for all workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results
	var errors []error
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}
