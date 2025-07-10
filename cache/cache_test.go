package cache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/cache"
	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)
	require.NotNil(t, cache)
}

func TestNewCache_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Should still work since context is reserved for future use
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)
	require.NotNil(t, cache)
}

func TestCache_BasicOperations(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test Set and Get
	ok := cache.Set("foo", "bar")
	require.True(t, ok)

	value, exists := cache.Get("foo")
	require.True(t, exists)
	require.Equal(t, "bar", value)

	// Test Delete
	cache.Delete("foo")

	_, exists = cache.Get("foo")
	require.False(t, exists)
}

func TestCache_NonExistentKey(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test getting non-existent key
	value, exists := cache.Get("non-existent")
	require.False(t, exists)
	require.Nil(t, value)
}

func TestCache_DeleteNonExistentKey(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Should not panic or error when deleting non-existent key
	require.NotPanics(t, func() {
		cache.Delete("non-existent")
	})
}

func TestCache_OverwriteValue(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Set initial value
	ok := cache.Set("key1", "value1")
	require.True(t, ok)

	// Overwrite with new value
	ok = cache.Set("key1", "value2")
	require.True(t, ok)

	// Verify new value
	value, exists := cache.Get("key1")
	require.True(t, exists)
	require.Equal(t, "value2", value)
}

func TestCache_DifferentValueTypes(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test different data types
	testCases := []struct {
		key   string
		value any
	}{
		{"string", "test string"},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"slice", []string{"a", "b", "c"}},
		{"map", map[string]int{"a": 1, "b": 2}},
		{"struct", struct{ Name string }{Name: "test"}},
		{"nil", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			ok := cache.Set(tc.key, tc.value)
			require.True(t, ok)

			value, exists := cache.Get(tc.key)
			require.True(t, exists)
			require.Equal(t, tc.value, value)
		})
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // readers, writers, deleters

	// Concurrent readers
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for range numOperations {
				key := "reader_key"
				cache.Get(key)
			}
		}(i)
	}

	// Concurrent writers
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for range numOperations {
				key := "writer_key"
				value := "value"
				cache.Set(key, value)
			}
		}(i)
	}

	// Concurrent deleters
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for range numOperations {
				key := "delete_key"
				cache.Delete(key)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - possible deadlock")
	}
}

func TestCache_TTLBehavior(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Set a value
	ok := cache.Set("ttl_test", "value")
	require.True(t, ok)

	// Should exist immediately
	value, exists := cache.Get("ttl_test")
	require.True(t, exists)
	require.Equal(t, "value", value)

	// Wait longer than TTL (1 minute + buffer)
	// Note: This test may be slow, but it verifies TTL functionality
	t.Run("TTL expiration", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping TTL test in short mode")
		}

		// Wait slightly longer than the 1-minute TTL
		time.Sleep(65 * time.Second)

		// Value should be expired
		_, exists := cache.Get("ttl_test")
		require.False(t, exists)
	})
}

func TestCache_MultipleKeysOperations(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Set multiple keys
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	values := []string{"value1", "value2", "value3", "value4", "value5"}

	for i, key := range keys {
		ok := cache.Set(key, values[i])
		require.True(t, ok)
	}

	// Verify all keys exist
	for i, key := range keys {
		value, exists := cache.Get(key)
		require.True(t, exists)
		require.Equal(t, values[i], value)
	}

	// Delete some keys
	cache.Delete("key2")
	cache.Delete("key4")

	// Verify deleted keys don't exist
	_, exists := cache.Get("key2")
	require.False(t, exists)
	_, exists = cache.Get("key4")
	require.False(t, exists)

	// Verify remaining keys still exist
	for _, key := range []string{"key1", "key3", "key5"} {
		_, exists := cache.Get(key)
		require.True(t, exists)
	}
}

func TestCache_LargeValues(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test with large string
	largeValue := make([]byte, 1024*1024) // 1MB
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}

	ok := cache.Set("large_value", largeValue)
	require.True(t, ok)

	value, exists := cache.Get("large_value")
	require.True(t, exists)
	require.Equal(t, largeValue, value)
}

func TestCache_EmptyStringKey(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Empty string key should work
	ok := cache.Set("", "empty key value")
	require.True(t, ok)

	value, exists := cache.Get("")
	require.True(t, exists)
	require.Equal(t, "empty key value", value)

	cache.Delete("")
	_, exists = cache.Get("")
	require.False(t, exists)
}

func TestCache_SetReturnValue(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test that Set always returns true (success)
	result := cache.Set("test_key", "test_value")
	require.True(t, result)

	// Test Set with nil value
	result = cache.Set("nil_key", nil)
	require.True(t, result)

	// Test Set with empty string value
	result = cache.Set("empty_key", "")
	require.True(t, result)
}

func TestCache_GetReturnValues(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Test Get with existing key
	cache.Set("existing_key", "value")
	value, exists := cache.Get("existing_key")
	require.True(t, exists)
	require.Equal(t, "value", value)

	// Test Get with non-existent key returns proper values
	value, exists = cache.Get("non_existent_key")
	require.False(t, exists)
	require.Nil(t, value)

	// Test Get after Delete
	cache.Delete("existing_key")
	value, exists = cache.Get("existing_key")
	require.False(t, exists)
	require.Nil(t, value)
}

func TestCache_DeleteIdempotency(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Set a value first
	cache.Set("delete_test", "value")

	// Verify it exists
	_, exists := cache.Get("delete_test")
	require.True(t, exists)

	// Delete it
	cache.Delete("delete_test")

	// Verify it's gone
	_, exists = cache.Get("delete_test")
	require.False(t, exists)

	// Delete again - should not panic or error
	require.NotPanics(t, func() {
		cache.Delete("delete_test")
	})

	// Multiple deletes of same key
	require.NotPanics(t, func() {
		cache.Delete("delete_test")
		cache.Delete("delete_test")
		cache.Delete("delete_test")
	})
}

func TestCache_SetOverwriteBehavior(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Set initial value
	result := cache.Set("overwrite_key", "original_value")
	require.True(t, result)

	// Verify initial value
	value, exists := cache.Get("overwrite_key")
	require.True(t, exists)
	require.Equal(t, "original_value", value)

	// Overwrite with different type
	result = cache.Set("overwrite_key", 42)
	require.True(t, result)

	// Verify overwritten value
	value, exists = cache.Get("overwrite_key")
	require.True(t, exists)
	require.Equal(t, 42, value)

	// Overwrite with nil
	result = cache.Set("overwrite_key", nil)
	require.True(t, result)

	// Verify nil value
	value, exists = cache.Get("overwrite_key")
	require.True(t, exists)
	require.Nil(t, value)
}

func TestCache_GetAfterSetSequence(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		key   string
		value any
	}{
		{"string_value", "str_key", "hello world"},
		{"integer_value", "int_key", 123},
		{"boolean_value", "bool_key", false},
		{"slice_value", "slice_key", []int{1, 2, 3}},
		{"nil_value", "nil_key", nil},
		{"empty_string", "empty_key", ""},
		{"zero_int", "zero_key", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the value
			result := cache.Set(tc.key, tc.value)
			require.True(t, result)

			// Immediately get the value
			retrievedValue, exists := cache.Get(tc.key)
			require.True(t, exists)
			require.Equal(t, tc.value, retrievedValue)
		})
	}
}

func TestCache_DeleteNonExistentKeys(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	// Delete keys that never existed
	nonExistentKeys := []string{
		"never_existed",
		"also_never_existed",
		"",
		"key_with_special_chars_!@#$%",
		"very_long_key_" + string(make([]byte, 1000)),
	}

	for _, key := range nonExistentKeys {
		require.NotPanics(t, func() {
			cache.Delete(key)
		})
	}
}

func TestCache_SetGetDeleteCycle(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(t, err)

	key := "cycle_test"
	values := []any{"first", 42, true, nil, []string{"a", "b"}}

	for i, value := range values {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			// Set
			result := cache.Set(key, value)
			require.True(t, result)

			// Get
			retrievedValue, exists := cache.Get(key)
			require.True(t, exists)
			require.Equal(t, value, retrievedValue)

			// Delete
			cache.Delete(key)

			// Verify deletion
			_, exists = cache.Get(key)
			require.False(t, exists)
		})
	}
}

// Benchmark tests
func BenchmarkCache_Set(b *testing.B) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(b, err)

	for b.Loop() {
		cache.Set("benchmark_key", "benchmark_value")
	}
}

func BenchmarkCache_Get(b *testing.B) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(b, err)

	// Pre-populate cache
	cache.Set("benchmark_key", "benchmark_value")

	for b.Loop() {
		cache.Get("benchmark_key")
	}
}

func BenchmarkCache_Delete(b *testing.B) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	require.NoError(b, err)

	for b.Loop() {
		cache.Set("benchmark_key", "benchmark_value")
		cache.Delete("benchmark_key")
	}
}
