# Toolshed Agent Guide

## Build/Test Commands
- `go test ./...` - Run all tests
- `go test ./package` - Run tests for specific package (e.g., `go test ./aes`)
- `go test -run TestName ./package` - Run single test (e.g., `go test -run TestGenerateAESKey_ValidLengths ./aes`)
- `go build ./...` - Build all packages
- `go mod tidy` - Clean up dependencies

## Architecture
- **Structure**: Go module with independent packages in separate directories
- **Packages**: aes, argon, base62, base64, clock, csv, hash, null, password, ulid
- **Testing**: Uses testify/require for assertions; test files follow `*_test.go` pattern
- **Dependencies**: Minimal external deps (oklog/ulid, wagslane/go-password-validator, golang.org/x/crypto)

## Code Style
- **Package naming**: lowercase, single words matching directory names
- **Imports**: Standard library first, then external packages
- **Error handling**: Use fmt.Errorf with %w verb for error wrapping
- **Comments**: Exported functions have descriptive comments starting with function name
- **Testing**: External test packages (package_test) with descriptive test function names
- **Constants**: Use const for default values (e.g., DefaultEntropy)
- **Base64**: Use base64.StdEncoding for consistency across packages
