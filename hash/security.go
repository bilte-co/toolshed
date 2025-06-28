package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// HMAC computes the HMAC of data using the specified key and algorithm.
func HMAC(data, key []byte, algorithm string) ([]byte, error) {
	var hashFunc func() hash.Hash
	
	algorithm = strings.ToLower(algorithm)
	switch algorithm {
	case "md5":
		hashFunc = md5.New
	case "sha1":
		hashFunc = sha1.New
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	case "blake2b":
		hashFunc = func() hash.Hash {
			h, _ := blake2b.New256(nil)
			return h
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}
	
	mac := hmac.New(hashFunc, key)
	mac.Write(data)
	return mac.Sum(nil), nil
}

// HMACWithOptions computes HMAC with custom output options.
func HMACWithOptions(data, key []byte, algorithm string, opts Options) (interface{}, error) {
	result, err := HMAC(data, key, algorithm)
	if err != nil {
		return nil, err
	}
	return formatOutput(result, "hmac-"+algorithm, opts)
}

// EqualConstantTime performs constant-time comparison of two byte slices.
// This prevents timing attacks when comparing sensitive data like hashes.
func EqualConstantTime(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// PasswordHashingOptions configures password hashing operations.
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

// DefaultPasswordOptions provides secure defaults for password hashing.
var DefaultPasswordOptions = PasswordHashingOptions{
	// PBKDF2 defaults
	PBKDF2Iterations: 100000,
	PBKDF2KeyLength:  32,
	PBKDF2SaltLength: 16,
	PBKDF2Algorithm:  "sha256",
	
	// scrypt defaults (recommended by RFC 7914)
	ScryptN:       32768, // CPU/memory cost parameter (2^15)
	ScryptR:       8,     // block size parameter
	ScryptP:       1,     // parallelization parameter
	ScryptKeyLen:  32,    // derived key length
	ScryptSaltLen: 16,    // salt length
	
	// bcrypt defaults
	BcryptCost: 12, // cost factor (2^12 iterations)
}

// PBKDF2Hash derives a key from a password using PBKDF2.
func PBKDF2Hash(password []byte, opts *PasswordHashingOptions) ([]byte, []byte, error) {
	if opts == nil {
		opts = &DefaultPasswordOptions
	}
	
	// Generate random salt
	salt := make([]byte, opts.PBKDF2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	
	return PBKDF2HashWithSalt(password, salt, opts)
}

// PBKDF2HashWithSalt derives a key from a password using PBKDF2 with the provided salt.
func PBKDF2HashWithSalt(password, salt []byte, opts *PasswordHashingOptions) ([]byte, []byte, error) {
	if opts == nil {
		opts = &DefaultPasswordOptions
	}
	
	var hashFunc func() hash.Hash
	algorithm := strings.ToLower(opts.PBKDF2Algorithm)
	switch algorithm {
	case "md5":
		hashFunc = md5.New
	case "sha1":
		hashFunc = sha1.New
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	case "blake2b":
		hashFunc = func() hash.Hash {
			h, _ := blake2b.New256(nil)
			return h
		}
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedAlgorithm, algorithm)
	}
	
	key := pbkdf2.Key(password, salt, opts.PBKDF2Iterations, opts.PBKDF2KeyLength, hashFunc)
	return key, salt, nil
}

// VerifyPBKDF2 verifies a password against a PBKDF2 hash.
func VerifyPBKDF2(password, salt, expectedHash []byte, opts *PasswordHashingOptions) bool {
	derivedKey, _, err := PBKDF2HashWithSalt(password, salt, opts)
	if err != nil {
		return false
	}
	return EqualConstantTime(derivedKey, expectedHash)
}

// ScryptHash derives a key from a password using scrypt.
func ScryptHash(password []byte, opts *PasswordHashingOptions) ([]byte, []byte, error) {
	if opts == nil {
		opts = &DefaultPasswordOptions
	}
	
	// Generate random salt
	salt := make([]byte, opts.ScryptSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	
	return ScryptHashWithSalt(password, salt, opts)
}

// ScryptHashWithSalt derives a key from a password using scrypt with the provided salt.
func ScryptHashWithSalt(password, salt []byte, opts *PasswordHashingOptions) ([]byte, []byte, error) {
	if opts == nil {
		opts = &DefaultPasswordOptions
	}
	
	key, err := scrypt.Key(password, salt, opts.ScryptN, opts.ScryptR, opts.ScryptP, opts.ScryptKeyLen)
	if err != nil {
		return nil, nil, fmt.Errorf("scrypt key derivation failed: %w", err)
	}
	
	return key, salt, nil
}

// VerifyScrypt verifies a password against a scrypt hash.
func VerifyScrypt(password, salt, expectedHash []byte, opts *PasswordHashingOptions) bool {
	derivedKey, _, err := ScryptHashWithSalt(password, salt, opts)
	if err != nil {
		return false
	}
	return EqualConstantTime(derivedKey, expectedHash)
}

// BcryptHash hashes a password using bcrypt.
func BcryptHash(password []byte, opts *PasswordHashingOptions) ([]byte, error) {
	cost := DefaultPasswordOptions.BcryptCost
	if opts != nil && opts.BcryptCost > 0 {
		cost = opts.BcryptCost
	}
	
	return bcrypt.GenerateFromPassword(password, cost)
}

// VerifyBcrypt verifies a password against a bcrypt hash.
func VerifyBcrypt(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

// GenerateSalt generates a cryptographically secure random salt.
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}
