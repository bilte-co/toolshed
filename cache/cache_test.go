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
