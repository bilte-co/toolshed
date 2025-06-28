// Package null provides utilities for converting sql.Null* types to pointers.
// These functions help bridge the gap between database NULL values and Go's pointer types,
// making it easier to work with optional fields in structs and JSON serialization.
//
// Example usage:
//
//	var ns sql.NullString
//	// ... populate ns from database query
//
//	// Convert to *string, nil if NULL
//	ptr := null.StringToPtr(ns)
//	if ptr != nil {
//		fmt.Println("Value:", *ptr)
//	} else {
//		fmt.Println("Value is NULL")
//	}
package null

import (
	"database/sql"
	"time"
)

// StringToPtr converts a sql.NullString to a *string.
// Returns a pointer to the string value if Valid is true, nil otherwise.
func StringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// Int64ToPtr converts a sql.NullInt64 to a *int64.
// Returns a pointer to the int64 value if Valid is true, nil otherwise.
func Int64ToPtr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

// Int32ToPtr converts a sql.NullInt32 to a *int32.
// Returns a pointer to the int32 value if Valid is true, nil otherwise.
func Int32ToPtr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

// Int16ToPtr converts a sql.NullInt16 to a *int16.
// Returns a pointer to the int16 value if Valid is true, nil otherwise.
func Int16ToPtr(ni sql.NullInt16) *int16 {
	if ni.Valid {
		return &ni.Int16
	}
	return nil
}

// Float64ToPtr converts a sql.NullFloat64 to a *float64.
// Returns a pointer to the float64 value if Valid is true, nil otherwise.
func Float64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

// Float32ToPtr converts a sql.NullFloat64 to a *float32 with type conversion.
// Returns a pointer to the float32 value if Valid is true, nil otherwise.
// Note: This performs a float64 to float32 conversion which may lose precision.
func Float32ToPtr(nf sql.NullFloat64) *float32 {
	if nf.Valid {
		val := float32(nf.Float64) // Explicit conversion
		return &val
	}
	return nil
}

// TimeToPtr converts a sql.NullTime to a *time.Time.
// Returns a pointer to the time.Time value if Valid is true, nil otherwise.
func TimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// BoolToPtr converts a sql.NullBool to a *bool.
// Returns a pointer to the bool value if Valid is true, nil otherwise.
func BoolToPtr(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}
