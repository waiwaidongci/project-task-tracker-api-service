package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	DatabasePath string
}

func Load() Config {
	return Config{
		Port:         envOrDefault("PORT", "18005"),
		DatabasePath: envOrDefault("DB_PATH", "todo.db"),
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
