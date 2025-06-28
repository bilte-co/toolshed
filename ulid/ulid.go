// Package ulid provides ULID (Universally Unique Lexicographically Sortable Identifier) generation
// and manipulation with base62 encoding. Unlike standard ULIDs which use base32 encoding,
// this package uses base62 for more compact string representation.
//
// ULIDs are sortable by creation time and contain both timestamp and random components,
// making them ideal for distributed systems where you need unique, sortable identifiers.
//
// Example usage:
//
//	// Create a new ULID with prefix
//	id, err := ulid.CreateULID("user", time.Now())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(id) // user_base62encodedULID
//
//	// Extract timestamp from ULID
//	timestamp, err := ulid.Timestamp(id)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Created at:", timestamp)
package ulid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/bilte-co/toolshed/base62"
	"github.com/oklog/ulid/v2"
)

// CreateULID generates a new ULID with an optional prefix and given timestamp.
// The ULID is encoded using base62 for compact representation. If prefix is not empty,
// it will be prepended to the ULID with an underscore separator.
func CreateULID(prefix string, timestamp time.Time) (string, error) {
	entropy := rand.Reader
	ms := ulid.Timestamp(timestamp)

	newUlid, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	// Encode raw 16 bytes as base62 instead of base32
	base62Encoded := base62.StdEncoding.Encode(newUlid[:])

	var sb strings.Builder
	if prefix != "" {
		sb.WriteString(prefix)
		sb.WriteString("_")
	}
	sb.WriteString(string(base62Encoded))

	return sb.String(), nil
}

// Decode converts a base62-encoded ULID string (with optional prefix) back to a ulid.ULID.
// The function automatically strips any prefix before decoding.
func Decode(s string) (ulid.ULID, error) {
	str, err := stripPrefix(s)
	if err != nil {
		return ulid.ULID{}, err
	}

	// Decode base62 string to raw ULID bytes
	decoded, err := base62.StdEncoding.DecodeString(str)
	if err != nil {
		return ulid.ULID{}, err
	}

	if len(decoded) != 16 {
		return ulid.ULID{}, errors.New("invalid ULID length")
	}

	return ulid.ULID(decoded), nil
}

// Timestamp extracts the timestamp component from a ULID string.
// Returns the time.Time when the ULID was created, automatically handling prefixes.
func Timestamp(id string) (time.Time, error) {
	str, err := stripPrefix(id)
	if err != nil {
		return time.Time{}, err
	}

	// Decode base62 string to raw ULID bytes
	decoded, err := base62.StdEncoding.DecodeString(str)
	if err != nil {
		return time.Time{}, err
	}

	if len(decoded) != 16 {
		return time.Time{}, errors.New("invalid ULID length")
	}

	// First 6 bytes are the timestamp (48-bit big-endian)
	timestampMs := binary.BigEndian.Uint64(append([]byte{0, 0}, decoded[:6]...))

	return time.UnixMilli(int64(timestampMs)), nil
}

// stripPrefix removes the prefix from a ULID string, if present.
// Returns the base62-encoded ULID portion without the prefix and underscore.
func stripPrefix(s string) (string, error) {
	var str string
	pos := strings.LastIndex(s, "_")
	if pos == -1 {
		str = s
	} else {
		str = s[pos+1:]
	}

	if str == "" {
		return "", errors.New("invalid ULID")
	}

	return str, nil
}
