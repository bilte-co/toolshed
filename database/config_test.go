package database_test

import (
	"os"
	"testing"
	"time"

	"github.com/bilte-co/toolshed/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_DatabaseConfig(t *testing.T) {
	config := &database.Config{
		Name: "testdb",
		User: "testuser",
	}

	result := config.DatabaseConfig()
	assert.Equal(t, config, result)
}

func TestNewConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *database.Config
	}{
		{
			name: "basic environment variables",
			envVars: map[string]string{
				"DB_NAME":     "testdb",
				"DB_USER":     "testuser",
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_SSLMODE":  "disable",
				"DB_PASSWORD": "testpass",
			},
			expected: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "localhost",
				Port:              "5432",
				SSLMode:           "disable",
				Password:          "testpass",
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "with SSL certificates",
			envVars: map[string]string{
				"DB_NAME":        "testdb",
				"DB_USER":        "testuser",
				"DB_SSLCERT":     "/path/to/cert.pem",
				"DB_SSLKEY":      "/path/to/key.pem",
				"DB_SSLROOTCERT": "/path/to/ca.pem",
			},
			expected: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				SSLCertPath:       "/path/to/cert.pem",
				SSLKeyPath:        "/path/to/key.pem",
				SSLRootCertPath:   "/path/to/ca.pem",
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "with pool settings",
			envVars: map[string]string{
				"DB_NAME":                     "testdb",
				"DB_POOL_MIN_CONNS":           "2",
				"DB_POOL_MAX_CONNS":           "10",
				"DB_POOL_MAX_CONN_LIFETIME":   "10m",
				"DB_POOL_MAX_CONN_IDLE_TIME":  "5m",
				"DB_POOL_HEALTH_CHECK_PERIOD": "30s",
			},
			expected: &database.Config{
				Name:               "testdb",
				PoolMinConnections: "2",
				PoolMaxConnections: "10",
				PoolMaxConnLife:    10 * time.Minute,
				PoolMaxConnIdle:    5 * time.Minute,
				PoolHealthCheck:    30 * time.Second,
				ConnectionTimeout:  0,
			},
		},
		{
			name: "with connection timeout",
			envVars: map[string]string{
				"DB_NAME":            "testdb",
				"DB_CONNECT_TIMEOUT": "30",
			},
			expected: &database.Config{
				Name:              "testdb",
				ConnectionTimeout: 30,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "invalid timeout defaults to 0",
			envVars: map[string]string{
				"DB_NAME":            "testdb",
				"DB_CONNECT_TIMEOUT": "invalid",
			},
			expected: &database.Config{
				Name:              "testdb",
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "invalid durations use defaults",
			envVars: map[string]string{
				"DB_NAME":                     "testdb",
				"DB_POOL_MAX_CONN_LIFETIME":   "invalid",
				"DB_POOL_MAX_CONN_IDLE_TIME":  "invalid",
				"DB_POOL_HEALTH_CHECK_PERIOD": "invalid",
			},
			expected: &database.Config{
				Name:              "testdb",
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearEnv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := database.NewConfigFromEnv()

			assert.Equal(t, tt.expected.Name, config.Name)
			assert.Equal(t, tt.expected.User, config.User)
			assert.Equal(t, tt.expected.Host, config.Host)
			assert.Equal(t, tt.expected.Port, config.Port)
			assert.Equal(t, tt.expected.SSLMode, config.SSLMode)
			assert.Equal(t, tt.expected.Password, config.Password)
			assert.Equal(t, tt.expected.SSLCertPath, config.SSLCertPath)
			assert.Equal(t, tt.expected.SSLKeyPath, config.SSLKeyPath)
			assert.Equal(t, tt.expected.SSLRootCertPath, config.SSLRootCertPath)
			assert.Equal(t, tt.expected.PoolMinConnections, config.PoolMinConnections)
			assert.Equal(t, tt.expected.PoolMaxConnections, config.PoolMaxConnections)
			assert.Equal(t, tt.expected.ConnectionTimeout, config.ConnectionTimeout)
			assert.Equal(t, tt.expected.PoolMaxConnLife, config.PoolMaxConnLife)
			assert.Equal(t, tt.expected.PoolMaxConnIdle, config.PoolMaxConnIdle)
			assert.Equal(t, tt.expected.PoolHealthCheck, config.PoolHealthCheck)

			// Cleanup
			clearEnv()
		})
	}
}

func TestNewConfigFromEnv_WithDSN(t *testing.T) {
	clearEnv()
	defer clearEnv()

	dsn := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	os.Setenv("DB_DSN", dsn)

	config := database.NewConfigFromEnv()

	assert.Equal(t, "testdb", config.Name)
	assert.Equal(t, "testuser", config.User)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "5432", config.Port)
	assert.Equal(t, "testpass", config.Password)
}

func TestNewConfigFromEnv_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *database.Config
	}{
		{
			name:    "completely empty environment",
			envVars: map[string]string{},
			expected: &database.Config{
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "negative timeout value",
			envVars: map[string]string{
				"DB_CONNECT_TIMEOUT": "-30",
			},
			expected: &database.Config{
				ConnectionTimeout: -30,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "zero timeout value",
			envVars: map[string]string{
				"DB_CONNECT_TIMEOUT": "0",
			},
			expected: &database.Config{
				ConnectionTimeout: 0,
				PoolMaxConnLife:   5 * time.Minute,
				PoolMaxConnIdle:   1 * time.Minute,
				PoolHealthCheck:   1 * time.Minute,
			},
		},
		{
			name: "invalid DSN in DB_DSN",
			envVars: map[string]string{
				"DB_DSN":  "invalid-dsn-format",
				"DB_NAME": "fallback-db", // These should be ignored when DSN is present
				"DB_USER": "fallback-user",
			},
			expected: nil, // Should return nil when DSN parsing fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv()
			defer clearEnv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := database.NewConfigFromEnv()

			if tt.expected == nil {
				assert.Nil(t, config)
			} else {
				require.NotNil(t, config)
				assert.Equal(t, tt.expected.Name, config.Name)
				assert.Equal(t, tt.expected.User, config.User)
				assert.Equal(t, tt.expected.Host, config.Host)
				assert.Equal(t, tt.expected.Port, config.Port)
				assert.Equal(t, tt.expected.Password, config.Password)
				assert.Equal(t, tt.expected.ConnectionTimeout, config.ConnectionTimeout)
				assert.Equal(t, tt.expected.PoolMaxConnLife, config.PoolMaxConnLife)
				assert.Equal(t, tt.expected.PoolMaxConnIdle, config.PoolMaxConnIdle)
				assert.Equal(t, tt.expected.PoolHealthCheck, config.PoolHealthCheck)
			}
		})
	}
}

func TestNewFromDSN(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		expected    *database.Config
		expectError bool
	}{
		{
			name: "valid PostgreSQL DSN",
			dsn:  "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
			expected: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name: "DSN without password",
			dsn:  "postgres://testuser@localhost:5432/testdb",
			expected: &database.Config{
				Name: "testdb",
				User: "testuser",
				Host: "localhost",
				Port: "5432",
			},
			expectError: false,
		},
		{
			name: "DSN with different port",
			dsn:  "postgres://testuser:testpass@localhost:3306/testdb",
			expected: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "3306",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name:        "invalid DSN",
			dsn:         "invalid-dsn",
			expected:    nil,
			expectError: true,
		},
		{
			name: "empty DSN",
			dsn:  "",
			expected: &database.Config{
				// pgx.ParseConfig with empty string falls back to environment/defaults
				Name:     "",
				User:     "inghamemerson", // System user as default
				Host:     "/private/tmp",  // Unix socket path
				Port:     "5432",
				Password: "",
			},
			expectError: false,
		},
		{
			name: "DSN without database name",
			dsn:  "postgres://testuser:testpass@localhost:5432/",
			expected: &database.Config{
				Name:     "",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name: "DSN with IPv6 host",
			dsn:  "postgres://testuser:testpass@[::1]:5432/testdb",
			expected: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "::1",
				Port:     "5432",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name: "DSN with special characters in password",
			dsn:  "postgres://testuser:p%40ss%21w%40rd@localhost:5432/testdb",
			expected: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "p@ss!w@rd",
			},
			expectError: false,
		},
		{
			name: "DSN with no user credentials",
			dsn:  "postgres://localhost:5432/testdb",
			expected: &database.Config{
				Name: "testdb",
				User: "inghamemerson", // Falls back to system user
				Host: "localhost",
				Port: "5432",
			},
			expectError: false,
		},
		{
			name: "DSN with default port",
			dsn:  "postgres://testuser:testpass@localhost/testdb",
			expected: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := database.NewFromDSN(tt.dsn)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				require.NoError(t, err)
				require.NotNil(t, config)
				assert.Equal(t, tt.expected.Name, config.Name)
				assert.Equal(t, tt.expected.User, config.User)
				assert.Equal(t, tt.expected.Host, config.Host)
				assert.Equal(t, tt.expected.Port, config.Port)
				assert.Equal(t, tt.expected.Password, config.Password)
			}
		})
	}
}

func TestConfig_ConnectionURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *database.Config
		expected string
	}{
		{
			name:     "nil config",
			config:   nil,
			expected: "",
		},
		{
			name: "basic configuration",
			config: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb",
		},
		{
			name: "with SSL mode",
			config: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
				SSLMode:  "require",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require",
		},
		{
			name: "with connection timeout",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "localhost",
				Port:              "5432",
				Password:          "testpass",
				ConnectionTimeout: 30,
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?connect_timeout=30",
		},
		{
			name: "with SSL certificates",
			config: &database.Config{
				Name:            "testdb",
				User:            "testuser",
				Host:            "localhost",
				Port:            "5432",
				Password:        "testpass",
				SSLCertPath:     "/path/to/cert.pem",
				SSLKeyPath:      "/path/to/key.pem",
				SSLRootCertPath: "/path/to/ca.pem",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslcert=%2Fpath%2Fto%2Fcert.pem&sslkey=%2Fpath%2Fto%2Fkey.pem&sslrootcert=%2Fpath%2Fto%2Fca.pem",
		},
		{
			name: "without port",
			config: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "localhost",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@localhost/testdb",
		},
		{
			name: "without user credentials",
			config: &database.Config{
				Name: "testdb",
				Host: "localhost",
				Port: "5432",
			},
			expected: "postgres://localhost:5432/testdb",
		},
		{
			name: "with all parameters",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "localhost",
				Port:              "5432",
				Password:          "testpass",
				SSLMode:           "require",
				ConnectionTimeout: 30,
				SSLCertPath:       "/path/to/cert.pem",
				SSLKeyPath:        "/path/to/key.pem",
				SSLRootCertPath:   "/path/to/ca.pem",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?connect_timeout=30&sslcert=%2Fpath%2Fto%2Fcert.pem&sslkey=%2Fpath%2Fto%2Fkey.pem&sslmode=require&sslrootcert=%2Fpath%2Fto%2Fca.pem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ConnectionURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_ConnectionURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		config   *database.Config
		expected string
	}{
		{
			name: "with empty database name",
			config: &database.Config{
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@localhost:5432",
		},
		{
			name: "with zero timeout",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "localhost",
				Port:              "5432",
				Password:          "testpass",
				ConnectionTimeout: 0,
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb",
		},
		{
			name: "with negative timeout",
			config: &database.Config{
				Name:              "testdb",
				User:              "testuser",
				Host:              "localhost",
				Port:              "5432",
				Password:          "testpass",
				ConnectionTimeout: -30,
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb",
		},
		{
			name: "with special characters in database name",
			config: &database.Config{
				Name:     "test-db_name",
				User:     "testuser",
				Host:     "localhost",
				Port:     "5432",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@localhost:5432/test-db_name",
		},
		{
			name: "with IPv6 host",
			config: &database.Config{
				Name:     "testdb",
				User:     "testuser",
				Host:     "::1",
				Port:     "5432",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@::1:5432/testdb",
		},
		{
			name: "with only user, no password",
			config: &database.Config{
				Name: "testdb",
				User: "testuser",
				Host: "localhost",
				Port: "5432",
			},
			expected: "postgres://testuser:@localhost:5432/testdb",
		},
		{
			name: "all empty fields",
			config: &database.Config{},
			expected: "postgres:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ConnectionURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func clearEnv() {
	envVars := []string{
		"DB_DSN",
		"DB_NAME",
		"DB_USER",
		"DB_HOST",
		"DB_PORT",
		"DB_SSLMODE",
		"DB_PASSWORD",
		"DB_SSLCERT",
		"DB_SSLKEY",
		"DB_SSLROOTCERT",
		"DB_POOL_MIN_CONNS",
		"DB_POOL_MAX_CONNS",
		"DB_CONNECT_TIMEOUT",
		"DB_POOL_MAX_CONN_LIFETIME",
		"DB_POOL_MAX_CONN_IDLE_TIME",
		"DB_POOL_HEALTH_CHECK_PERIOD",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
