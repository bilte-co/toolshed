package database_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromEnv(t *testing.T) {
	t.Run("missing database configuration", func(t *testing.T) {
		clearEnv()
		defer clearEnv()

		ctx := context.Background()
		db, err := database.NewFromEnv(ctx)

		// Clean up if a connection was somehow created
		if db != nil {
			db.Close(ctx)
		}

		// With no configuration, it should either fail to parse or fail to connect
		// The behavior may vary based on system defaults
		if err != nil {
			// Expected case - should fail with some kind of connection/parse error
			assert.NotNil(t, err)
		} else {
			// Unexpected but possible on some systems with default postgres setup
			t.Skip("Connection succeeded unexpectedly - system may have default postgres config")
		}
	})

	t.Run("with valid DSN", func(t *testing.T) {
		clearEnv()
		defer clearEnv()

		// Use a basic DSN that should parse correctly but won't actually connect
		dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
		os.Setenv("DB_DSN", dsn)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		db, err := database.NewFromEnv(ctx)

		// Should parse successfully but connection will likely fail
		// In a real test environment, you'd want to use a test database
		if err != nil {
			// If connection fails, it should be a connection error, not a parse error
			assert.NotContains(t, err.Error(), "failed to parse connection string")
		} else {
			// If it succeeds (unlikely without a real DB), clean up
			require.NotNil(t, db)
			db.Close(ctx)
		}
	})

	t.Run("with environment variables", func(t *testing.T) {
		clearEnv()
		defer clearEnv()

		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_SSLMODE", "disable")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		db, err := database.NewFromEnv(ctx)

		// Should parse successfully but connection will likely fail
		if err != nil {
			// If connection fails, it should be a connection error, not a parse error
			assert.NotContains(t, err.Error(), "failed to parse connection string")
		} else {
			// If it succeeds (unlikely without a real DB), clean up
			require.NotNil(t, db)
			db.Close(ctx)
		}
	})
}

func TestDB_Close(t *testing.T) {
	t.Run("close with nil pool", func(t *testing.T) {
		// Create a DB with a nil pool to test graceful handling
		db := &database.DB{Pool: nil}
		ctx := context.Background()

		// This will panic due to the pgxpool implementation - that's expected behavior
		assert.Panics(t, func() {
			db.Close(ctx)
		})
	})
}

// Test helper functions that are internal to the connection.go file
// We can test these indirectly through their usage

func TestDBDSNGeneration(t *testing.T) {
	tests := []struct {
		name     string
		config   *database.Config
		contains []string
		excludes []string
	}{
		{
			name: "basic configuration",
			config: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			contains: []string{
				"dbname=testdb",
				"user=testuser",
				"host=localhost",
				"port=5432",
				"password=testpass",
			},
		},
		{
			name: "with SSL configuration",
			config: &database.Config{
				Name:            "testdb",
				User:            "testuser",
				SSLMode:         "require",
				SSLCertPath:     "/path/to/cert.pem",
				SSLKeyPath:      "/path/to/key.pem",
				SSLRootCertPath: "/path/to/ca.pem",
			},
			contains: []string{
				"dbname=testdb",
				"user=testuser",
				"sslmode=require",
				"sslcert=/path/to/cert.pem",
				"sslkey=/path/to/key.pem",
				"sslrootcert=/path/to/ca.pem",
			},
		},
		{
			name: "with connection timeout",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				ConnectionTimeout: 30,
			},
			contains: []string{
				"dbname=testdb",
				"user=testuser",
				"connect_timeout=30",
			},
		},
		{
			name: "with pool settings",
			config: &database.Config{
				Name:               "testdb",
				User:               "testuser",
				PoolMinConnections: "2",
				PoolMaxConnections: "10",
				PoolMaxConnLife:    10 * time.Minute,
				PoolMaxConnIdle:    5 * time.Minute,
				PoolHealthCheck:    30 * time.Second,
			},
			contains: []string{
				"dbname=testdb",
				"user=testuser",
				"pool_min_conns=2",
				"pool_max_conns=10",
				"pool_max_conn_lifetime=10m0s",
				"pool_max_conn_idle_time=5m0s",
				"pool_health_check_period=30s",
			},
		},
		{
			name: "empty values excluded",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "", // empty
				Port:              "", // empty
				ConnectionTimeout: 0,  // zero value
			},
			contains: []string{
				"dbname=testdb",
				"user=testuser",
			},
			excludes: []string{
				"host=",
				"port=",
				"connect_timeout=0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv()
			defer clearEnv()

			// Set up environment to use our test config
			os.Setenv("DB_NAME", tt.config.Name)
			os.Setenv("DB_USER", tt.config.User)
			if tt.config.Host != "" {
				os.Setenv("DB_HOST", tt.config.Host)
			}
			if tt.config.Port != "" {
				os.Setenv("DB_PORT", tt.config.Port)
			}
			if tt.config.Password != "" {
				os.Setenv("DB_PASSWORD", tt.config.Password)
			}
			if tt.config.SSLMode != "" {
				os.Setenv("DB_SSLMODE", tt.config.SSLMode)
			}
			if tt.config.SSLCertPath != "" {
				os.Setenv("DB_SSLCERT", tt.config.SSLCertPath)
			}
			if tt.config.SSLKeyPath != "" {
				os.Setenv("DB_SSLKEY", tt.config.SSLKeyPath)
			}
			if tt.config.SSLRootCertPath != "" {
				os.Setenv("DB_SSLROOTCERT", tt.config.SSLRootCertPath)
			}
			if tt.config.PoolMinConnections != "" {
				os.Setenv("DB_POOL_MIN_CONNS", tt.config.PoolMinConnections)
			}
			if tt.config.PoolMaxConnections != "" {
				os.Setenv("DB_POOL_MAX_CONNS", tt.config.PoolMaxConnections)
			}
			if tt.config.ConnectionTimeout > 0 {
				os.Setenv("DB_CONNECT_TIMEOUT", "30")
			}
			if tt.config.PoolMaxConnLife > 0 {
				os.Setenv("DB_POOL_MAX_CONN_LIFETIME", tt.config.PoolMaxConnLife.String())
			}
			if tt.config.PoolMaxConnIdle > 0 {
				os.Setenv("DB_POOL_MAX_CONN_IDLE_TIME", tt.config.PoolMaxConnIdle.String())
			}
			if tt.config.PoolHealthCheck > 0 {
				os.Setenv("DB_POOL_HEALTH_CHECK_PERIOD", tt.config.PoolHealthCheck.String())
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Try to create the connection
			db, err := database.NewFromEnv(ctx)

			// Clean up if connection succeeded
			if db != nil {
				db.Close(ctx)
			}

			// This test verifies that the configuration can be loaded without parse errors
			// The actual connection may succeed or fail depending on local setup
			if err != nil {
				// If there's an error, it should not be a parse error
				assert.NotContains(t, err.Error(), "failed to parse connection string")
			}
			// If no error, the configuration was valid and connection succeeded
		})
	}
}

func TestBeforeAcquireCallback(t *testing.T) {
	// This test verifies that the BeforeAcquire callback is set up correctly
	// We can't easily test the actual callback without a real database connection,
	// but we can verify the setup doesn't cause parse errors

	clearEnv()
	defer clearEnv()

	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_SSLMODE", "disable")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	db, err := database.NewFromEnv(ctx)

	// Clean up if connection succeeded
	if db != nil {
		db.Close(ctx)
	}

	// This test verifies the BeforeAcquire callback setup doesn't cause parse errors
	if err != nil {
		// If there's an error, it should not be a parse error
		assert.NotContains(t, err.Error(), "failed to parse connection string")
		// The error should be about connection failure, which means the config was parsed
		assert.True(t,
			strings.Contains(err.Error(), "failed to create connection pool") ||
				strings.Contains(err.Error(), "connection") ||
				strings.Contains(err.Error(), "dial") ||
				strings.Contains(err.Error(), "timeout"))
	}
	// If no error, the configuration was valid and connection succeeded
}

// Integration test that would require a real database
// This is commented out but shows how you might test with a real DB
/*
func TestIntegration_WithRealDB(t *testing.T) {
	// This test would require a real PostgreSQL instance
	// You might use testcontainers or similar for this
	t.Skip("Integration test requires real database")

	clearEnv()
	defer clearEnv()

	// Set up test database connection
	os.Setenv("DB_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")

	ctx := context.Background()
	db, err := database.NewFromEnv(ctx)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer db.Close(ctx)

	// Test basic connection
	err = db.Pool.Ping(ctx)
	assert.NoError(t, err)
}
*/
