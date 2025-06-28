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
