[![Go Reference](https://pkg.go.dev/badge/github.com/bilte-co/toolshed.svg)](https://pkg.go.dev/github.com/bilte-co/toolshed)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilte-co/toolshed)](https://goreportcard.com/report/github.com/bilte-co/toolshed)
[![codecov](https://codecov.io/gh/bilte-co/toolshed/branch/main/graph/badge.svg?token=U9JB17FAGT)](https://codecov.io/gh/bilte-co/toolshed)

# Toolshed CLI

A command-line interface for the Toolshed Go utility library, featuring robust hashing, ULID generation, password validation, and encryption operations with enterprise-grade reliability and performance.

## Features

- **Multiple Hash Algorithms**: SHA-256, SHA-512, SHA-1, MD5, BLAKE2b
- **ULID Generation**: Sortable, time-based unique identifiers with custom prefixes
- **AES Encryption**: Secure file encryption/decryption with AES-GCM
- **Flexible Input Sources**: Strings, files, directories, stdin
- **HMAC Support**: Secure message authentication codes
- **Hash Validation**: Verify file integrity against expected checksums
- **Password Strength Validation**: Entropy-based password security checking
- **Constant-Time Comparison**: Secure hash comparison preventing timing attacks
- **Production Ready**: Structured logging, error handling, concurrency safety
- **User-Friendly**: Colorized output, progress spinners, clear error messages

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/bilte-co/toolshed.git
cd toolshed

# Build and install
make install
```

### Development Build

```bash
make dev
```

This will format, vet, and build the binary to `build/toolshed`.

## Usage

### Basic Commands

```bash
# Hash a string
toolshed hash string "Hello, World!" --algo sha256

# Hash a file
toolshed hash file document.pdf --algo sha512

# Hash from stdin
echo "Hello" | toolshed hash file -

# Hash a directory recursively
toolshed hash dir /path/to/directory --recursive

# Compute HMAC
toolshed hash hmac "sensitive data" --key "secret-key" --algo sha256

# Validate file integrity
toolshed hash validate document.pdf --expected a1b2c3d4... --algo sha256

# Compare two hashes securely
toolshed hash compare a1b2c3d4... e5f6a7b8...

# Check password strength
toolshed password check "MySecurePassword123!"

# Check password with custom entropy requirement
toolshed password check "password" --entropy 50

# Generate ULIDs
toolshed ulid create
toolshed ulid create --prefix "user"
toolshed ulid create --timestamp "2023-01-01T00:00:00Z"

# Extract timestamp from ULID
toolshed ulid timestamp "user_30KMu42XfVhcsuTE9VgFm"
toolshed ulid timestamp "user_30KMu42XfVhcsuTE9VgFm" --format unix
```

### Advanced Options

```bash
# Different output formats
toolshed hash string "test" --format base64
toolshed hash string "test" --format hex --prefix

# Verbose logging
toolshed --verbose hash file large-file.zip

# Version information
toolshed --version
```

### Output Formats

- `hex` (default): Hexadecimal encoding
- `base64`: Base64 encoding
- `raw`: Raw bytes (binary output)

Use `--prefix` to include the algorithm name in output (e.g., `sha256:a1b2c3...`).

## Examples

### File Integrity Verification

```bash
# Generate checksum
toolshed hash file important-file.zip --algo sha256 > checksum.txt

# Later, verify integrity
toolshed hash validate important-file.zip --expected $(cat checksum.txt) --algo sha256
```

### Secure Data Authentication

```bash
# Create HMAC for API authentication
toolshed hash hmac "user=alice&action=transfer&amount=100" \
  --key "api-secret-key" --algo sha256
```

### Directory Monitoring

```bash
# Create directory snapshot
toolshed hash dir /etc/config --recursive > config-baseline.hash

# Later, check for changes
CURRENT=$(toolshed hash dir /etc/config --recursive)
BASELINE=$(cat config-baseline.hash)
toolshed hash compare "$CURRENT" "$BASELINE"
```

### Password Security Validation

```bash
# Check password from command line
toolshed password check "MySecurePassword123!"

# Check password from stdin (secure input)
echo "SuperSecretPassword!" | toolshed password check

# Check with custom entropy requirement
toolshed password check "password123" --entropy 70

# Interactive password checking
toolshed password check
# Prompts: Enter password to check:

# Check multiple passwords from file
cat passwords.txt | while read -r pwd; do
  echo "$pwd" | toolshed password check --entropy 65
done
```

### ULID Operations

```bash
# Generate a new ULID
toolshed ulid create

# Generate ULID with prefix
toolshed ulid create --prefix "user"

# Generate ULID with custom timestamp
toolshed ulid create --timestamp "2023-01-01T00:00:00Z" --prefix "order"

# Extract timestamp from ULID (RFC3339 format)
toolshed ulid timestamp "user_30KMu42XfVhcsuTE9VgFm"

# Extract timestamp in Unix format
toolshed ulid timestamp "user_30KMu42XfVhcsuTE9VgFm" --format unix

# Extract timestamp in Unix milliseconds
toolshed ulid timestamp "user_30KMu42XfVhcsuTE9VgFm" --format unixmilli

# Extract timestamp from stdin
echo "user_30KMu42XfVhcsuTE9VgFm" | toolshed ulid timestamp -

# Batch processing ULIDs
cat ulids.txt | while read -r ulid; do
  echo "ULID: $ulid, Created: $(echo "$ulid" | toolshed ulid timestamp -)"
done
```

## Security Features

- **Constant-Time Comparison**: Prevents timing attacks when comparing hashes
- **Input Sanitization**: Safe path handling and validation
- **Memory Efficient**: Streaming for large files, no full-file loading
- **Secure Defaults**: SHA-256 default algorithm, secure HMAC implementation

## Development

### Building

```bash
# Development build with checks
make dev

# Production release (multiple platforms)
make release

# Run tests with coverage
make test

# Format and lint
make fmt lint
```

### Project Structure

```
toolshed/
├── main.go              # CLI entry point
├── internal/cli/        # CLI command implementations
│   ├── aes.go           # AES encryption commands
│   ├── context.go       # Shared context
│   ├── hash.go          # Hash commands
│   ├── password.go      # Password commands
│   ├── serve.go         # File server commands
│   ├── ulid.go          # ULID commands
│   └── version.go       # Version handling
├── hash/                # Hash utility package
├── password/            # Password utility package
├── ulid/                # ULID utility package
├── aes/                 # AES encryption package
├── Makefile             # Build automation
└── README.md
```

## Dependencies

- [kong](https://github.com/alecthomas/kong) - Command-line parsing
- [tint](https://github.com/lmittmann/tint) - Colored structured logging
- [spinner](https://github.com/briandowns/spinner) - Progress indicators

## Performance

The CLI is optimized for production use:

- **Large Files**: Streaming I/O, configurable buffer sizes
- **Directory Hashing**: Concurrent file processing
- **Memory Efficient**: Constant memory usage regardless of file size
- **Fast Algorithms**: Hardware-accelerated when available

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make dev` to ensure quality
6. Submit a pull request

## License

This project is part of the Toolshed utility library. See the main repository for license information.
