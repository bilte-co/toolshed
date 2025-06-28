package ulid

import (
	"strings"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/base62"
	"github.com/stretchr/testify/require"
)

func TestCreateULID_NoPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("", ts)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	ulidBytes, err := base62.StdEncoding.DecodeString(id)
	require.NoError(t, err)
	require.Len(t, ulidBytes, 16)
}

func TestCreateULID_WithPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("test", ts)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(id, "test_"))

	parts := strings.SplitN(id, "_", 2)
	require.Len(t, parts, 2)
	ulidBytes, err := base62.StdEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	require.Len(t, ulidBytes, 16)
}

func TestDecode_ValidULID(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("myobj", ts)
	require.NoError(t, err)

	ulidVal, err := Decode(id)
	require.NoError(t, err)
	require.Equal(t, 16, len(ulidVal))
}

func TestDecode_NoPrefix(t *testing.T) {
	ts := time.Now()
	id, err := CreateULID("", ts)
	require.NoError(t, err)

	ulidVal, err := Decode(id)
	require.NoError(t, err)
	require.Equal(t, 16, len(ulidVal))
}

func TestDecode_InvalidBase62(t *testing.T) {
	_, err := Decode("prefix_***invalid***")
	require.Error(t, err)
}

func TestDecode_InvalidLength(t *testing.T) {
	_, err := Decode("prefix_" + base62.StdEncoding.EncodeToString([]byte("short")))
	require.ErrorContains(t, err, "invalid ULID length")
}

func TestTimestamp_Extract(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	id, err := CreateULID("ts", now)
	require.NoError(t, err)

	extracted, err := Timestamp(id)
	require.NoError(t, err)

	// Should match to the millisecond
	require.WithinDuration(t, now, extracted, time.Millisecond)
}

func TestTimestamp_InvalidID(t *testing.T) {
	_, err := Timestamp("bad_string")
	require.Error(t, err)
}

func Test_stripPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"abc_def_ghi", "ghi", false},
		{"abc", "abc", false},
		{"no_suffix_", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		out, err := stripPrefix(tt.input)
		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, out)
		}
	}
}
