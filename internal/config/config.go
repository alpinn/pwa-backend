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

type JWTConfig struct {
    Secret   string
    KeyID    string
    Issuer   string
    Audience string
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

func NewJWTConfig(secret string) *JWTConfig {
    return &JWTConfig{
        Secret:   secret,
        KeyID:    "powersync-key",
        Issuer:   "pwa-backend",
        Audience: "https://6940ebf14011d65924582a54.powersync.journeyapps.com",
    }
}