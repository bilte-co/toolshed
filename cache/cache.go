// Package cache provides high-performance in-memory caching with TTL support.
// It offers both a simple interface and optimized implementations using the otter cache library
// for production workloads requiring fast access times and automatic expiration.
//
// Example usage:
//
//	// Create a new cache instance
//	ctx := context.Background()
//	cache, err := cache.NewCache(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Store a value
//	cache.Set("user:123", userObject)
//
//	// Retrieve a value
//	if value, exists := cache.Get("user:123"); exists {
//		user := value.(User)
//		fmt.Printf("Found user: %+v\n", user)
//	}
//
//	// Remove a value
//	cache.Delete("user:123")
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/maypok86/otter"
)

// Cache defines the interface for cache implementations.
// It provides basic cache operations for storing, retrieving, and removing key-value pairs.
// All implementations should be safe for concurrent use.
type Cache interface {
	// Get retrieves a value by key from the cache.
	// Returns the value and true if the key exists, or nil and false if not found.
	Get(key string) (any, bool)

	// Set stores a value with the specified key in the cache.
	// Returns true if the operation was successful, false otherwise.
	Set(key string, value any) bool

	// Delete removes a key-value pair from the cache.
	// No error is returned if the key doesn't exist.
	Delete(key string)
}

// InMemoryCache is a simple thread-safe in-memory cache implementation.
// It uses a map with read-write mutex for concurrent access protection.
// This implementation does not support TTL or automatic eviction.
type InMemoryCache struct {
	mu   sync.RWMutex   // Mutex for protecting concurrent access
	data map[string]any // Internal storage for key-value pairs
}

// NewCache creates a new high-performance cache instance using the otter library.
// The cache is configured with a maximum capacity of 1,000 entries, 1-minute TTL,
// and statistics collection enabled. All entries have equal cost (1) for eviction purposes.
// The context parameter is reserved for future use and cancellation support.
// Returns an error if cache initialization fails, though this is unlikely with current configuration.
func NewCache(ctx context.Context) (Cache, error) {
	cache, err := otter.MustBuilder[string, any](1_000).
		CollectStats().
		Cost(func(key string, value any) uint32 {
			return 1
		}).
		WithTTL(time.Minute).
		Build()
	if err != nil {
		panic(err)
	}

	return cache, nil
}

// Get retrieves a value from the in-memory cache by key.
// It uses a read lock to allow concurrent reads while preventing data races.
// Returns the stored value and true if the key exists, or nil and false if not found.
func (c *InMemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

// Set stores a key-value pair in the in-memory cache.
// It uses a write lock to ensure thread safety during modifications.
// Returns true to indicate the operation was successful.
func (c *InMemoryCache) Set(key string, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	return true
}

// Delete removes a key-value pair from the in-memory cache.
// It uses a write lock to ensure thread safety during modifications.
// No error is returned if the key doesn't exist.
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
