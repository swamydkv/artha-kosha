package config

import (
	"os"
	"testing"
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
}
