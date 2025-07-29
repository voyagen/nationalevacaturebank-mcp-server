package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration values
type Config struct {
	BaseURL       string
	Timeout       time.Duration
	LogLevel      string
	MaxRetries    int
	ServerName    string
	ServerVersion string
}

// Load loads configuration from environment variables or defaults
func Load() (*Config, error) {
	cfg := &Config{
		BaseURL:       getEnv("NVB_BASE_URL", "https://api.nationalevacaturebank.nl"),
		Timeout:       getDuration("NVB_TIMEOUT", 30*time.Second),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		MaxRetries:    getInt("MAX_RETRIES", 3),
		ServerName:    getEnv("SERVER_NAME", "Nationale Vacature Bank"),
		ServerVersion: getEnv("SERVER_VERSION", "1.0.0"),
	}

	return cfg, cfg.validate()
}

// validate ensures the configuration is valid
func (c *Config) validate() error {
	// Add validation logic here if needed
	// For now, we accept all values
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDuration gets a duration from environment or returns default
func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getInt gets an integer from environment or returns default
func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}