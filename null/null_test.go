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
	t.Run("valid string", func(t *testing.T) {
		str := "hello"
		result := null.StringToPtr(sql.NullString{String: str, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, str, *result)
	})

	t.Run("empty string", func(t *testing.T) {
		str := ""
		result := null.StringToPtr(sql.NullString{String: str, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, str, *result)
	})

	t.Run("null string", func(t *testing.T) {
		result := null.StringToPtr(sql.NullString{Valid: false})
		require.Nil(t, result)
	})

	t.Run("null with non-empty string field", func(t *testing.T) {
		// Even if String field is set, should return nil when Valid=false
		result := null.StringToPtr(sql.NullString{String: "ignored", Valid: false})
		require.Nil(t, result)
	})
}

func TestInt64ToPtr(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		val := int64(123)
		result := null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("negative value", func(t *testing.T) {
		val := int64(-456)
		result := null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("zero value", func(t *testing.T) {
		val := int64(0)
		result := null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("max int64", func(t *testing.T) {
		val := int64(9223372036854775807) // math.MaxInt64
		result := null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("min int64", func(t *testing.T) {
		val := int64(-9223372036854775808) // math.MinInt64
		result := null.Int64ToPtr(sql.NullInt64{Int64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("null value", func(t *testing.T) {
		result := null.Int64ToPtr(sql.NullInt64{Valid: false})
		require.Nil(t, result)
	})

	t.Run("null with non-zero int64 field", func(t *testing.T) {
		result := null.Int64ToPtr(sql.NullInt64{Int64: 999, Valid: false})
		require.Nil(t, result)
	})
}

func TestInt32ToPtr(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		val := int32(456)
		result := null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("negative value", func(t *testing.T) {
		val := int32(-789)
		result := null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("zero value", func(t *testing.T) {
		val := int32(0)
		result := null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("max int32", func(t *testing.T) {
		val := int32(2147483647) // math.MaxInt32
		result := null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("min int32", func(t *testing.T) {
		val := int32(-2147483648) // math.MinInt32
		result := null.Int32ToPtr(sql.NullInt32{Int32: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("null value", func(t *testing.T) {
		result := null.Int32ToPtr(sql.NullInt32{Valid: false})
		require.Nil(t, result)
	})
}

func TestInt16ToPtr(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		val := int16(789)
		result := null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("negative value", func(t *testing.T) {
		val := int16(-123)
		result := null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("zero value", func(t *testing.T) {
		val := int16(0)
		result := null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("max int16", func(t *testing.T) {
		val := int16(32767) // math.MaxInt16
		result := null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("min int16", func(t *testing.T) {
		val := int16(-32768) // math.MinInt16
		result := null.Int16ToPtr(sql.NullInt16{Int16: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("null value", func(t *testing.T) {
		result := null.Int16ToPtr(sql.NullInt16{Valid: false})
		require.Nil(t, result)
	})
}

func TestFloat64ToPtr(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		val := 123.456
		result := null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("negative value", func(t *testing.T) {
		val := -789.123
		result := null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("zero value", func(t *testing.T) {
		val := 0.0
		result := null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("very small value", func(t *testing.T) {
		val := 1e-100
		result := null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("very large value", func(t *testing.T) {
		val := 1e100
		result := null.Float64ToPtr(sql.NullFloat64{Float64: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("null value", func(t *testing.T) {
		result := null.Float64ToPtr(sql.NullFloat64{Valid: false})
		require.Nil(t, result)
	})
}

func TestFloat32ToPtr(t *testing.T) {
	t.Run("positive value", func(t *testing.T) {
		val := float32(123.456)
		result := null.Float32ToPtr(sql.NullFloat64{Float64: float64(val), Valid: true})
		require.NotNil(t, result)
		require.InDelta(t, val, *result, 1e-6)
	})

	t.Run("negative value", func(t *testing.T) {
		val := float32(-789.123)
		result := null.Float32ToPtr(sql.NullFloat64{Float64: float64(val), Valid: true})
		require.NotNil(t, result)
		require.InDelta(t, val, *result, 1e-6)
	})

	t.Run("zero value", func(t *testing.T) {
		val := float32(0.0)
		result := null.Float32ToPtr(sql.NullFloat64{Float64: float64(val), Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("precision loss", func(t *testing.T) {
		// Test value that loses precision in float64->float32 conversion
		original := 123.456789123456789 // High precision float64
		result := null.Float32ToPtr(sql.NullFloat64{Float64: original, Valid: true})
		require.NotNil(t, result)
		// Should be approximately equal to float32 conversion
		expected := float32(original)
		require.Equal(t, expected, *result)
	})

	t.Run("large value", func(t *testing.T) {
		val := float32(1e30)
		result := null.Float32ToPtr(sql.NullFloat64{Float64: float64(val), Valid: true})
		require.NotNil(t, result)
		require.InDelta(t, val, *result, 1e24) // Allow for precision differences
	})

	t.Run("null value", func(t *testing.T) {
		result := null.Float32ToPtr(sql.NullFloat64{Valid: false})
		require.Nil(t, result)
	})
}

func TestTimeToPtr(t *testing.T) {
	t.Run("current time", func(t *testing.T) {
		now := time.Now()
		result := null.TimeToPtr(sql.NullTime{Time: now, Valid: true})
		require.NotNil(t, result)
		require.WithinDuration(t, now, *result, time.Millisecond)
	})

	t.Run("unix epoch", func(t *testing.T) {
		epoch := time.Unix(0, 0)
		result := null.TimeToPtr(sql.NullTime{Time: epoch, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, epoch, *result)
	})

	t.Run("zero time", func(t *testing.T) {
		zero := time.Time{}
		result := null.TimeToPtr(sql.NullTime{Time: zero, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, zero, *result)
	})

	t.Run("future time", func(t *testing.T) {
		future := time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
		result := null.TimeToPtr(sql.NullTime{Time: future, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, future, *result)
	})

	t.Run("past time", func(t *testing.T) {
		past := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		result := null.TimeToPtr(sql.NullTime{Time: past, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, past, *result)
	})

	t.Run("null time", func(t *testing.T) {
		result := null.TimeToPtr(sql.NullTime{Valid: false})
		require.Nil(t, result)
	})
}

func TestBoolToPtr(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		val := true
		result := null.BoolToPtr(sql.NullBool{Bool: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("false value", func(t *testing.T) {
		val := false
		result := null.BoolToPtr(sql.NullBool{Bool: val, Valid: true})
		require.NotNil(t, result)
		require.Equal(t, val, *result)
	})

	t.Run("null value", func(t *testing.T) {
		result := null.BoolToPtr(sql.NullBool{Valid: false})
		require.Nil(t, result)
	})

	t.Run("null with true bool field", func(t *testing.T) {
		// Even if Bool field is true, should return nil when Valid=false
		result := null.BoolToPtr(sql.NullBool{Bool: true, Valid: false})
		require.Nil(t, result)
	})
}
