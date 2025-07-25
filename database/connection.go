package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bilte-co/toolshed/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents a database connection with connection pooling.
// It wraps a pgxpool.Pool to provide high-performance PostgreSQL connectivity
// with automatic connection management and health checking.
type DB struct {
	Pool *pgxpool.Pool // PostgreSQL connection pool
}

// NewFromEnv creates a new database connection using environment configuration.
// It automatically configures connection pooling, health checks, and connection validation.
// The context is used for connection establishment and should have appropriate timeout.
// Returns an error if the configuration is invalid or connection fails.
func NewFromEnv(ctx context.Context) (*DB, error) {
	cfg := NewConfigFromEnv()

	pgxConfig, err := pgxpool.ParseConfig(dbDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// BeforeAcquire is called before before a connection is acquired from the
	// pool. It must return true to allow the acquisition or false to indicate that
	// the connection should be destroyed and a different connection should be
	// acquired.
	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Ping the connection to see if it is still valid. Ping returns an error if
		// it fails.
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close gracefully closes the database connection pool.
// It logs the closure and ensures all connections are properly released.
// The context can be used to set a timeout for the close operation.
func (db *DB) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("🔌 Closing connection pool.")
	db.Pool.Close()
}

// dbDSN converts a Config to a PostgreSQL DSN string.
// It formats all configuration parameters into a space-separated key=value format
// suitable for pgx driver consumption.
func dbDSN(cfg *Config) string {
	vals := dbValues(cfg)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

// setIfNotEmpty adds a key-value pair to the map only if the value is not empty.
// This helper prevents empty configuration values from being included in the DSN.
func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

// setIfPositive adds a key-value pair to the map only if the integer value is positive.
// This helper prevents zero or negative values from being included in the DSN.
func setIfPositive(m map[string]string, key string, val int) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

// setIfPositiveDuration adds a key-value pair to the map only if the duration is positive.
// This helper prevents zero or negative durations from being included in the DSN.
func setIfPositiveDuration(m map[string]string, key string, d time.Duration) {
	if d > 0 {
		m[key] = d.String()
	}
}

// dbValues converts a Config struct to a map of PostgreSQL connection parameters.
// It maps the struct fields to their corresponding pgx parameter names and
// only includes non-empty/positive values to avoid configuration conflicts.
func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", cfg.Host)
	setIfNotEmpty(p, "port", cfg.Port)
	setIfNotEmpty(p, "sslmode", cfg.SSLMode)
	setIfPositive(p, "connect_timeout", cfg.ConnectionTimeout)
	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslcert", cfg.SSLCertPath)
	setIfNotEmpty(p, "sslkey", cfg.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", cfg.SSLRootCertPath)
	setIfNotEmpty(p, "pool_min_conns", cfg.PoolMinConnections)
	setIfNotEmpty(p, "pool_max_conns", cfg.PoolMaxConnections)
	setIfPositiveDuration(p, "pool_max_conn_lifetime", cfg.PoolMaxConnLife)
	setIfPositiveDuration(p, "pool_max_conn_idle_time", cfg.PoolMaxConnIdle)
	setIfPositiveDuration(p, "pool_health_check_period", cfg.PoolHealthCheck)
	return p
}
