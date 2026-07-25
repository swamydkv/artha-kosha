package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Test default values
	os.Unsetenv("PORT")
	os.Unsetenv("FINANCE_API_ALLOWED_ORIGINS")

	c := Load()
	if c.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", c.Port)
	}
	if c.AllowedOrigins != "" {
		t.Errorf("expected default empty allowed origins, got %s", c.AllowedOrigins)
	}

	// Test overridden values
	os.Setenv("PORT", "9090")
	os.Setenv("FINANCE_API_ALLOWED_ORIGINS", "http://localhost:3000")

	c = Load()
	if c.Port != "9090" {
		t.Errorf("expected port 9090, got %s", c.Port)
	}
	if c.AllowedOrigins != "http://localhost:3000" {
		t.Errorf("expected origins http://localhost:3000, got %s", c.AllowedOrigins)
	}

	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("FINANCE_API_ALLOWED_ORIGINS")

	// Test getEnvInt
	os.Setenv("INT_VAR", "42")
	if v := getEnvInt("INT_VAR", 10); v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
	os.Setenv("INT_VAR", "invalid")
	if v := getEnvInt("INT_VAR", 10); v != 10 {
		t.Errorf("expected 10, got %d", v)
	}

	// Test getEnvDuration
	os.Setenv("DUR_VAR", "10s")
	if v := getEnvDuration("DUR_VAR", time.Second); v != 10*time.Second {
		t.Errorf("expected 10s, got %v", v)
	}
	os.Setenv("DUR_VAR", "invalid")
	if v := getEnvDuration("DUR_VAR", time.Second); v != time.Second {
		t.Errorf("expected 1s, got %v", v)
	}
}
