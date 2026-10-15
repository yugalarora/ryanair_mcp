package config

import (
	"os"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	os.Unsetenv("RYANAIR_CURRENCY")
	os.Unsetenv("DEFAULT_CURRENCY")
	os.Unsetenv("RYANAIR_RADIUS_KM")
	os.Unsetenv("RYANAIR_MAX_OPTIONS")

	cfg := LoadConfig()
	if cfg.DefaultCurrency != "EUR" {
		t.Errorf("expected default currency EUR, got %s", cfg.DefaultCurrency)
	}
	if cfg.DefaultRadiusKm != 150.0 {
		t.Errorf("expected default radius 150.0, got %f", cfg.DefaultRadiusKm)
	}
	if cfg.DefaultMaxOptions != 15 {
		t.Errorf("expected default max options 15, got %d", cfg.DefaultMaxOptions)
	}
}

func TestLoadConfigCustomEnv(t *testing.T) {
	os.Setenv("RYANAIR_CURRENCY", "GBP")
	os.Setenv("RYANAIR_RADIUS_KM", "200")
	os.Setenv("RYANAIR_MAX_OPTIONS", "25")
	defer func() {
		os.Unsetenv("RYANAIR_CURRENCY")
		os.Unsetenv("RYANAIR_RADIUS_KM")
		os.Unsetenv("RYANAIR_MAX_OPTIONS")
	}()

	cfg := LoadConfig()
	if cfg.DefaultCurrency != "GBP" {
		t.Errorf("expected currency GBP, got %s", cfg.DefaultCurrency)
	}
	if cfg.DefaultRadiusKm != 200.0 {
		t.Errorf("expected radius 200.0, got %f", cfg.DefaultRadiusKm)
	}
	if cfg.DefaultMaxOptions != 25 {
		t.Errorf("expected max options 25, got %d", cfg.DefaultMaxOptions)
	}
}
