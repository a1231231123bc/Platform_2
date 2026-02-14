package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRES_IN", "24h")
	os.Setenv("PORT", "8080")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN")
		os.Unsetenv("PORT")
	}()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "postgresql://test:test@localhost:5432/test", cfg.DatabaseURL)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, 24*time.Hour, cfg.JWTExpiresIn)
	assert.Equal(t, "8080", cfg.Port)
}

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRES_IN")
	os.Unsetenv("PORT")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "", cfg.DatabaseURL)
	assert.Equal(t, "", cfg.JWTSecret)
	assert.Equal(t, 168*time.Hour, cfg.JWTExpiresIn)
	assert.Equal(t, "3000", cfg.Port)
}
