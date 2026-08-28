package config

import "os"

// Config contains process-level settings loaded from environment variables.
type Config struct {
	Address      string
	DatabaseURL  string
	Environment  string
	ResourcesDir string
	WebDir       string
}

func Load() Config {
	return Config{
		Address:      envOrDefault("SERVER_ADDRESS", ":8080"),
		DatabaseURL:  envOrDefault("DATABASE_URL", "postgres://workbench:workbench@localhost:5432/workbench?sslmode=disable"),
		Environment:  envOrDefault("APP_ENV", "development"),
		ResourcesDir: envOrDefault("RESOURCES_DIR", "resources"),
		WebDir:       envOrDefault("WEB_DIR", "web"),
	}
}

func envOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
