package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheck_WeakPassword(t *testing.T) {
	ok, msg := Check("123456")
	require.False(t, ok)
	require.NotEmpty(t, msg)
	require.Contains(t, msg, "insecure password")
}

func TestCheck_StrongPassword(t *testing.T) {
	ok, msg := Check("G@7e*vS93^8!bdT2")
	require.True(t, ok)
	require.Empty(t, msg)
}

func TestCheckEntropy_HighThreshold(t *testing.T) {
	ok, msg := CheckEntropy("strongPassword123!", 120.0)
	require.False(t, ok)
	require.NotEmpty(t, msg)
}

func TestCheckEntropy_LowThreshold(t *testing.T) {
	ok, msg := CheckEntropy("weakbutpassable", 20.0)
	require.True(t, ok)
	require.Empty(t, msg)
}

func TestCheck_EmptyPassword(t *testing.T) {
	ok, msg := Check("")
	require.False(t, ok)
	require.NotEmpty(t, msg)
}

func TestCheck_RepetitivePassword(t *testing.T) {
	ok, msg := Check("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	require.False(t, ok)
	require.Contains(t, msg, "more special characters")
}
