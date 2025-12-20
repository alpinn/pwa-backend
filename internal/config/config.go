package config

import (
	"os"
	"time"
)

type Config struct {
    DatabaseURL        string
    JWTSecret          string
    Port               string
    DBConnectTimeout   time.Duration
    DBMaxRetries       int
}

func Load() *Config {
    return &Config{
        DatabaseURL:      getEnv("DATABASE_URL", ""),
        JWTSecret:        getEnv("JWT_SECRET", "biskuat"),
        Port:             getEnv("PORT", "8080"),
        DBConnectTimeout: 30 * time.Second,
        DBMaxRetries:     5,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}