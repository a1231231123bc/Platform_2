package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
	Port         string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "168h"))
	if err != nil {
		expiresIn = 168 * time.Hour
	}

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		JWTExpiresIn: expiresIn,
		Port:         getEnv("PORT", "3000"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
