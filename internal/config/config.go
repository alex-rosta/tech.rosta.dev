package config

import (
	"os"
	"strings"
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

func Load() *Config {
	allowedUsers := []string{}
	if users := os.Getenv("ALLOWED_USERS"); users != "" {
		allowedUsers = strings.Split(users, ",")
		for i, user := range allowedUsers {
			allowedUsers[i] = strings.TrimSpace(user)
		}
	}

	return &Config{
		PostsDir:       getEnv("POSTS_DIR", "posts"),
		GitHubClientID: getEnv("CLIENT_ID", ""),
		GitHubSecret:   getEnv("CLIENT_SECRET", ""),
		SessionSecret:  getEnv("SESSION_SECRET", "your-secret-key-change-in-production"),
		AllowedUsers:   allowedUsers,
		BaseURL:        getEnv("BASE_URL", "http://localhost:3000"),
		Port:           getEnv("PORT", "3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
