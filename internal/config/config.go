package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL             string
	MaxConnections  int32
	MinConnections  int32
	MaxConnLifetime time.Duration
}

type Option func(*Config)

func LoadConfig() (*Config, error) {
	// Load .env file if exists
	if err := godotenv.Load(".env"); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
		log.Println("No .env file found, using environment variables")
	}

	// Create config with default options
	cfg := New(
		WithServerPort(getEnv("SERVER_PORT", ":8080")),
		WithDatabaseURL(getEnv("DB_URL", "")),
		WithMaxConnections(parseInt32(getEnv("DB_MAX_CONNECTIONS", "10"))),
		WithMinConnections(parseInt32(getEnv("DB_MIN_CONNECTIONS", "2"))),
		WithMaxConnLifetime(parseDuration(getEnv("DB_MAX_CONN_LIFETIME", "1h"))),
	)

	return cfg, nil
}

func New(options ...Option) *Config {
	cfg := &Config{}

	for _, option := range options {
		option(cfg)
	}

	return cfg
}

func WithServerPort(port string) Option {
	return func(c *Config) {
		c.Server.Port = port
	}
}

func WithDatabaseURL(url string) Option {
	return func(c *Config) {
		c.Database.URL = url
	}
}

func WithMaxConnections(max int32) Option {
	return func(c *Config) {
		c.Database.MaxConnections = max
	}
}

func WithMinConnections(min int32) Option {
	return func(c *Config) {
		c.Database.MinConnections = min
	}
}

func WithMaxConnLifetime(duration time.Duration) Option {
	return func(c *Config) {
		c.Database.MaxConnLifetime = duration
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseInt32(s string) int32 {
	var i int32
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Hour
	}
	return d
}
