// Package database provides PostgreSQL database configuration and connection management.
// It supports environment-based configuration, DSN parsing, and connection pooling
// using the pgx driver for optimal performance and reliability.
//
// Example usage:
//
//	// Create database connection from environment variables
//	ctx := context.Background()
//	db, err := database.NewFromEnv(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close(ctx)
//
//	// Create configuration from environment
//	config := database.NewConfigFromEnv()
//	connectionURL := config.ConnectionURL()
//
//	// Create configuration from DSN
//	config, err := database.NewFromDSN("postgres://user:pass@localhost/mydb")
//	if err != nil {
//		log.Fatal(err)
//	}
package database

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bilte-co/toolshed/logging"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Config holds database connection configuration parameters.
// It supports PostgreSQL connection settings including SSL, connection pooling,
// and timeout configurations for production database deployments.
type Config struct {
	Name               string        // Database name
	User               string        // Database user
	Host               string        // Database host
	Port               string        // Database port
	SSLMode            string        // SSL mode (disable, require, verify-ca, verify-full)
	ConnectionTimeout  int           // Connection timeout in seconds
	Password           string        // Database password
	SSLCertPath        string        // Path to SSL certificate file
	SSLKeyPath         string        // Path to SSL key file
	SSLRootCertPath    string        // Path to SSL root certificate file
	PoolMinConnections string        // Minimum connections in pool
	PoolMaxConnections string        // Maximum connections in pool
	PoolMaxConnLife    time.Duration // Maximum connection lifetime
	PoolMaxConnIdle    time.Duration // Maximum connection idle time
	PoolHealthCheck    time.Duration // Health check period for connections
}

// DatabaseConfig returns the database configuration.
// This method provides a consistent interface for accessing configuration settings.
func (c *Config) DatabaseConfig() *Config {
	return c
}

// NewConfigFromEnv creates a new database configuration from environment variables.
// It loads configuration from .env file if present and supports both individual
// environment variables and a complete DSN via DB_DSN.
// If DB_DSN is provided, it takes precedence over individual variables.
// Default values are applied for connection pool settings when not specified.
func NewConfigFromEnv() *Config {
	logger := logging.NewLoggerFromEnv()
	config := &Config{}

	err := godotenv.Load()
	if err != nil {
		logger.Warn("🤯 failed to load environment variables", "error", err)
	}

	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		config, _ := NewFromDSN(dsn)
		return config
	}

	config.Name = os.Getenv("DB_NAME")
	config.User = os.Getenv("DB_USER")
	config.Host = os.Getenv("DB_HOST")
	config.Port = os.Getenv("DB_PORT")
	config.SSLMode = os.Getenv("DB_SSLMODE")
	config.Password = os.Getenv("DB_PASSWORD")
	config.SSLCertPath = os.Getenv("DB_SSLCERT")
	config.SSLKeyPath = os.Getenv("DB_SSLKEY")
	config.SSLRootCertPath = os.Getenv("DB_SSLROOTCERT")
	config.PoolMinConnections = os.Getenv("DB_POOL_MIN_CONNS")
	config.PoolMaxConnections = os.Getenv("DB_POOL_MAX_CONNS")

	timeout, err := strconv.Atoi(os.Getenv("DB_CONNECT_TIMEOUT"))
	if err != nil {
		config.ConnectionTimeout = 0
	} else {
		config.ConnectionTimeout = timeout
	}

	poolMaxConnLife, err := time.ParseDuration(os.Getenv("DB_POOL_MAX_CONN_LIFETIME"))
	if err != nil {
		config.PoolMaxConnLife = 5 * time.Minute
	} else {
		config.PoolMaxConnLife = poolMaxConnLife
	}

	poolMaxConnIdle, err := time.ParseDuration(os.Getenv("DB_POOL_MAX_CONN_IDLE_TIME"))
	if err != nil {
		config.PoolMaxConnIdle = 1 * time.Minute
	} else {
		config.PoolMaxConnIdle = poolMaxConnIdle
	}

	poolHealthCheck, err := time.ParseDuration(os.Getenv("DB_POOL_HEALTH_CHECK_PERIOD"))
	if err != nil {
		config.PoolHealthCheck = 1 * time.Minute
	} else {
		config.PoolHealthCheck = poolHealthCheck
	}

	return config
}

// NewFromDSN creates a new database configuration from a PostgreSQL DSN (Data Source Name).
// The DSN should be in the format: postgres://user:password@host:port/database?options
// Returns an error if the DSN cannot be parsed by the pgx driver.
func NewFromDSN(dsn string) (*Config, error) {
	cfg := &Config{}

	vals, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	port := vals.Config.Port
	// convert to string
	portStr := strconv.FormatUint(uint64(port), 10)

	cfg.Name = vals.Config.Database
	cfg.User = vals.Config.User
	cfg.Host = vals.Config.Host
	cfg.Password = vals.Config.Password
	cfg.Port = portStr

	return cfg, nil
}

// ConnectionURL generates a PostgreSQL connection URL from the configuration.
// Returns a properly formatted postgres:// URL with all configured parameters.
// Returns an empty string if the configuration is nil.
func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	q := u.Query()
	if v := c.ConnectionTimeout; v > 0 {
		q.Add("connect_timeout", strconv.Itoa(v))
	}
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}
	if v := c.SSLCertPath; v != "" {
		q.Add("sslcert", v)
	}
	if v := c.SSLKeyPath; v != "" {
		q.Add("sslkey", v)
	}
	if v := c.SSLRootCertPath; v != "" {
		q.Add("sslrootcert", v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
