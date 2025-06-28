# HashUtil - Production-Grade Go Hashing Library

A comprehensive, flexible, and extensible Go library for hashing operations with security features and parallel processing capabilities.

## Features

### Core Hashing Functions

- **HashString**: Hash strings using various algorithms
- **HashBytes**: Hash byte slices
- **HashFile**: Hash individual files
- **HashReader**: Hash data from io.Reader streams
- **HashDir**: Hash directory contents deterministically

### Supported Algorithms

- SHA-1 ⚠️ (with security warning)
- SHA-256
- SHA-512
- BLAKE2b
- MD5 ⚠️ (with security warning)
- Custom algorithms via registration

### Security Features

- **HMAC**: Hash-based Message Authentication Code
- **Constant-time comparison**: Prevents timing attacks
- **Password hashing**:
  - PBKDF2 with configurable parameters
  - scrypt with secure defaults
  - bcrypt with configurable cost

### Advanced Features

- **Multiple output formats**: Raw bytes, hex, base64
- **Algorithm prefixing**: Optional algorithm prefix in output
- **Parallel processing**: Batch file hashing with worker pools
- **Archive support**: Hash compressed archives (.zip, .tar.gz, .tar)
- **Checksum validation**: Validate files against known checksums
- **Incremental hashing**: Hasher wrapper for streaming operations

### Extensibility

- **Custom algorithms**: Register your own hash implementations
- **Configurable options**: Customize behavior per operation

## Installation

```bash
go get github.com/bilte-co/toolshed/hash
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/bilte-co/toolshed/hash"
)

func main() {
    // Basic string hashing
    result, err := hash.HashString("hello world", "sha256")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hash: %x\n", result)

    // Hash with custom options
    formatted, err := hash.HashStringWithOptions("hello", "sha256", hash.Options{
        Format: hash.FormatHex,
        Prefix: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Formatted: %s\n", formatted)
}
```

## Examples

### Basic Hashing

```go
// Hash a string
hash, err := hash.HashString("hello world", "sha256")

// Hash bytes
hash, err := hash.HashBytes([]byte("data"), "sha512")

// Hash a file
hash, err := hash.HashFile("/path/to/file.txt", "blake2b")

// Hash a directory
hash, err := hash.HashDir("/path/to/dir", "sha256", true) // recursive
```

### Security Features

```go
// HMAC
mac, err := hash.HMAC(data, key, "sha256")

// Constant-time comparison
if hash.EqualConstantTime(hash1, hash2) {
    fmt.Println("Hashes match")
}

// Password hashing with PBKDF2
key, salt, err := hash.PBKDF2Hash(password, nil)
if hash.VerifyPBKDF2(password, salt, key, nil) {
    fmt.Println("Password verified")
}

// Bcrypt
hashed, err := hash.BcryptHash(password, &hash.PasswordHashingOptions{
    BcryptCost: 12,
})
if hash.VerifyBcrypt(password, hashed) {
    fmt.Println("Password verified")
}
```

### Incremental Hashing

```go
hasher, err := hash.NewHasher("sha256")
if err != nil {
    log.Fatal(err)
}

hasher.Write([]byte("hello"))
hasher.Write([]byte(" "))
hasher.Write([]byte("world"))

result := hasher.SumHex() // Get hex-encoded result
```

### Parallel File Processing

```go
filePaths := []string{"file1.txt", "file2.txt", "file3.txt"}

// Hash files in parallel with 4 workers
result := hash.HashFilesInParallel(filePaths, "sha256", 4)

for _, fileResult := range result.Results {
    if fileResult.Error == nil {
        fmt.Printf("File: %s, Hash: %x\n", fileResult.Path, fileResult.Hash)
    }
}
```

### File Validation

```go
// Validate a file against a known checksum
err := hash.ValidateFileChecksum("file.txt", "expected-hash-here", "sha256")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}

// Validate multiple files in parallel
checksums := []hash.FileChecksum{
    {Path: "file1.txt", ExpectedHash: "hash1"},
    {Path: "file2.txt", ExpectedHash: "hash2"},
}
errors := hash.ValidateFilesInParallel(checksums, "sha256", 4)
```

### Custom Hash Algorithms

```go
// Register a custom hasher
hash.RegisterHasher("custom-algo", func() hash.Hash {
    return sha256.New() // Your custom implementation here
})

// Use the custom hasher
result, err := hash.HashString("data", "custom-algo")
```

### Archive Support

```go
// Hash archive contents
hash, err := hash.HashArchive("archive.zip", "sha256")
hash, err := hash.HashArchive("archive.tar.gz", "sha256")

// Hash compressed file
hash, err := hash.HashCompressedFile("file.gz", "sha256")
```

## Configuration Options

### Options Structure

```go
type Options struct {
    Format     Format // FormatRaw, FormatHex, FormatBase64
    Prefix     bool   // Add algorithm prefix to output
    Workers    int    // Number of parallel workers
    BufferSize int    // Buffer size for I/O operations
}
```

### Password Hashing Options

```go
type PasswordHashingOptions struct {
    // PBKDF2 options
    PBKDF2Iterations int
    PBKDF2KeyLength  int
    PBKDF2SaltLength int
    PBKDF2Algorithm  string

    // scrypt options
    ScryptN        int
    ScryptR        int
    ScryptP        int
    ScryptKeyLen   int
    ScryptSaltLen  int

    // bcrypt options
    BcryptCost int
}
```

## Security Considerations

### Algorithm Security

- **MD5** and **SHA-1** are considered cryptographically broken and should only be used for non-security purposes
- **SHA-256** and **SHA-512** are recommended for general use
- **BLAKE2b** provides excellent performance and security
- The library automatically logs warnings when insecure algorithms are used

### Password Hashing Best Practices

- **bcrypt**: Simple to use, automatically handles salting
- **PBKDF2**: Widely supported, configurable iterations
- **scrypt**: Memory-hard function, resistant to hardware attacks

### Timing Attack Protection

- Use `EqualConstantTime()` for comparing sensitive data like hashes or MACs
- Never use standard byte slice comparison for cryptographic purposes

## Performance

The library is optimized for performance with:

- Parallel file processing with configurable worker pools
- Efficient I/O with configurable buffer sizes
- Streaming support for large data sets
- Memory-efficient archive processing

## Error Handling

All functions return detailed errors for:

- Unsupported algorithms
- File I/O errors
- Invalid input parameters
- Cryptographic failures

## Testing

Run the comprehensive test suite:

```bash
go test -v ./...
go test -race ./...
go test -bench=. ./...
```

## Dependencies

The library uses minimal external dependencies:

- `golang.org/x/crypto` for BLAKE2b, scrypt, bcrypt, and PBKDF2
- Go standard library for core functionality

## License

This library is part of the aviation.dev project and follows the same licensing terms.

## Contributing

1. Ensure all tests pass
2. Add tests for new functionality
3. Follow Go best practices
4. Update documentation for new features
5. Consider security implications of changes

## Security Reporting

For security-related issues, please follow the aviation.dev security reporting guidelines.
