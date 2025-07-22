package config

import (
	"os"
)

type Config struct {
	PostsDir       string
	GitHubClientID string
	GitHubSecret   string
	SessionSecret  string
	AllowedUsers   []string
	BaseURL        string
	Port           string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadEnv() *Config {
	return &Config{
		PostsDir: getEnv("POSTS_DIR", "posts"),
		Port:     getEnv("PORT", "3000"),
	}
}
