package hash_test

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"log"
	"os"
	"path/filepath"

	hashutil "github.com/bilte-co/toolshed/hash"
)

func ExampleHashString() {
	// Basic string hashing
	result, err := hashutil.HashString("hello world", "sha256")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SHA-256: %x\n", result)

	// Output:
	// SHA-256: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
}

func ExampleHashStringWithOptions() {
	// Hash with custom formatting and prefix
	result, err := hashutil.HashStringWithOptions("hello", "sha256", hashutil.Options{
		Format: hashutil.FormatHex,
		Prefix: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s\n", result)

	// Base64 format
	result, err = hashutil.HashStringWithOptions("hello", "sha256", hashutil.Options{
		Format: hashutil.FormatBase64,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base64: %s\n", result)
}

func ExampleNewHasher() {
	// Create a hasher for incremental hashing
	hasher, err := hashutil.NewHasher("sha256")
	if err != nil {
		log.Fatal(err)
	}

	// Write data incrementally
	hasher.Write([]byte("hello"))
	hasher.Write([]byte(" "))
	hasher.Write([]byte("world"))

	// Get result
	result := hasher.SumHex()
	fmt.Printf("Incremental hash: %s\n", result)
	fmt.Printf("Algorithm: %s\n", hasher.Algorithm())
}

func ExampleHMAC() {
	data := []byte("important message")
	key := []byte("secret-key")

	// Compute HMAC
	mac, err := hashutil.HMAC(data, key, "sha256")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HMAC: %x\n", mac)

	// Verify HMAC
	verification, err := hashutil.HMAC(data, key, "sha256")
	if err != nil {
		log.Fatal(err)
	}

	if hashutil.EqualConstantTime(mac, verification) {
		fmt.Println("HMAC verification: SUCCESS")
	}
}

func ExamplePBKDF2Hash() {
	password := []byte("user-password")

	// Hash password with PBKDF2
	key, salt, err := hashutil.PBKDF2Hash(password, nil) // Use default options
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Salt length: %d bytes\n", len(salt))
	fmt.Printf("Key length: %d bytes\n", len(key))

	// Verify password
	if hashutil.VerifyPBKDF2(password, salt, key, nil) {
		fmt.Println("Password verification: SUCCESS")
	}
}

func ExampleBcryptHash() {
	password := []byte("user-password")

	// Hash password with bcrypt
	hashed, err := hashutil.BcryptHash(password, &hashutil.PasswordHashingOptions{
		BcryptCost: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Bcrypt hash length: %d bytes\n", len(hashed))

	// Verify password
	if hashutil.VerifyBcrypt(password, hashed) {
		fmt.Println("Bcrypt verification: SUCCESS")
	}
}

func ExampleHashFilesInParallel() {
	// Create temporary files for demonstration
	tmpDir, err := os.MkdirTemp("", "hash-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	var filePaths []string

	for i, filename := range files {
		path := filepath.Join(tmpDir, filename)
		content := fmt.Sprintf("Content of file %d", i+1)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			log.Fatal(err)
		}
		filePaths = append(filePaths, path)
	}

	// Hash files in parallel
	result := hashutil.HashFilesInParallel(filePaths, "sha256", 2)

	fmt.Printf("Processed %d files\n", len(result.Results))
	fmt.Printf("Errors: %d\n", len(result.Errors))

	for _, fileResult := range result.Results {
		if fileResult.Error == nil {
			fmt.Printf("File: %s, Hash: %x\n",
				filepath.Base(fileResult.Path),
				fileResult.Hash[:8]) // Show first 8 bytes
		}
	}
}

func ExampleValidateFileChecksum() {
	// Create a temporary file
	tmpDir, err := os.MkdirTemp("", "hash-validate")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		log.Fatal(err)
	}

	// Get the file's hash
	fileHash, err := hashutil.HashFile(testFile, "sha256")
	if err != nil {
		log.Fatal(err)
	}

	expectedHash := fmt.Sprintf("%x", fileHash)

	// Validate checksum
	if err := hashutil.ValidateFileChecksum(testFile, expectedHash, "sha256"); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("File checksum validation: SUCCESS")
	}
}

func ExampleRegisterHasher() {
	// Register a custom hasher (for demonstration, we'll use an existing one)
	hashutil.RegisterHasher("custom-sha256", func() hash.Hash {
		return sha256.New()
	})

	// Use the custom hasher
	result, err := hashutil.HashString("test", "custom-sha256")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Custom hasher result: %x\n", result)
}

func ExampleHashDir() {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "hash-dir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files and subdirectories
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(tmpDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			log.Fatal(err)
		}
	}

	// Hash directory non-recursively
	hash1, err := hashutil.HashDir(tmpDir, "sha256", false)
	if err != nil {
		log.Fatal(err)
	}

	// Hash directory recursively
	hash2, err := hashutil.HashDir(tmpDir, "sha256", true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Non-recursive hash: %x\n", hash1[:8])
	fmt.Printf("Recursive hash: %x\n", hash2[:8])
	fmt.Printf("Hashes are different: %t\n", !hashutil.EqualConstantTime(hash1, hash2))
}
