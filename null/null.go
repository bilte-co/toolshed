package null

import (
	"database/sql"
	"time"
)

func StringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func Int64ToPtr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func Int32ToPtr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

func Int16ToPtr(ni sql.NullInt16) *int16 {
	if ni.Valid {
		return &ni.Int16
	}
	return nil
}

func Float64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func Float32ToPtr(nf sql.NullFloat64) *float32 {
	if nf.Valid {
		val := float32(nf.Float64) // Explicit conversion
		return &val
	}
	return nil
}

func TimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func BoolToPtr(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}
