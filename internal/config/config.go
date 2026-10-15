package config

import (
	"os"
	"strconv"
)

// Config stores global runtime parameters for the MCP server.
type Config struct {
	DefaultCurrency   string
	DefaultOrigin     string
	DefaultRadiusKm   float64
	DefaultMaxOptions int
	Transport         string
	Port              string
}

// LoadConfig reads configuration from environment variables with sensible defaults.
func LoadConfig() Config {
	cfg := Config{
		DefaultCurrency:   getEnv("RYANAIR_CURRENCY", getEnv("DEFAULT_CURRENCY", "EUR")),
		DefaultOrigin:     getEnv("RYANAIR_DEFAULT_ORIGIN", getEnv("DEFAULT_ORIGIN", "")),
		DefaultRadiusKm:   getEnvFloat("RYANAIR_RADIUS_KM", getEnvFloat("DEFAULT_RADIUS_KM", 150.0)),
		DefaultMaxOptions: getEnvInt("RYANAIR_MAX_OPTIONS", getEnvInt("DEFAULT_MAX_OPTIONS", 15)),
		Transport:         getEnv("MCP_TRANSPORT", "stdio"),
		Port:              getEnv("MCP_PORT", "8080"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil && f > 0 {
			return f
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}
