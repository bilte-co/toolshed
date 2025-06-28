package null_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/null"
	"github.com/stretchr/testify/require"
)

// ptr is a helper to create a pointer to a value.
func ptr[T any](v T) *T {
	return &v
}

func TestStringToPtr(t *testing.T) {
	str := "hello"
	require.Equal(t, ptr(str), null.StringToPtr(sql.NullString{String: str, Valid: true}))
	require.Nil(t, null.StringToPtr(sql.NullString{Valid: false}))
}

func TestInt64ToPtr(t *testing.T) {
	val := int64(123)
	require.Equal(t, ptr(val), null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true}))
	require.Nil(t, null.Int64ToPtr(sql.NullInt64{Valid: false}))
}

func TestInt32ToPtr(t *testing.T) {
	val := int32(456)
	require.Equal(t, ptr(val), null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true}))
	require.Nil(t, null.Int32ToPtr(sql.NullInt32{Valid: false}))
}

func TestInt16ToPtr(t *testing.T) {
	val := int16(789)
	require.Equal(t, ptr(val), null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true}))
	require.Nil(t, null.Int16ToPtr(sql.NullInt16{Valid: false}))
}

func TestFloat64ToPtr(t *testing.T) {
	val := 123.456
	require.Equal(t, ptr(val), null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true}))
	require.Nil(t, null.Float64ToPtr(sql.NullFloat64{Valid: false}))
}

func TestFloat32ToPtr(t *testing.T) {
	val := float32(123.456)
	result := null.Float32ToPtr(sql.NullFloat64{Float64: float64(val), Valid: true})
	require.NotNil(t, result)
	require.InDelta(t, val, *result, 1e-6)

	require.Nil(t, null.Float32ToPtr(sql.NullFloat64{Valid: false}))
}

func TestTimeToPtr(t *testing.T) {
	now := time.Now()
	result := null.TimeToPtr(sql.NullTime{Time: now, Valid: true})
	require.NotNil(t, result)
	require.WithinDuration(t, now, *result, time.Millisecond)

	require.Nil(t, null.TimeToPtr(sql.NullTime{Valid: false}))
}

func TestBoolToPtr(t *testing.T) {
	val := true
	require.Equal(t, ptr(val), null.BoolToPtr(sql.NullBool{Bool: val, Valid: true}))
	require.Nil(t, null.BoolToPtr(sql.NullBool{Valid: false}))
}
