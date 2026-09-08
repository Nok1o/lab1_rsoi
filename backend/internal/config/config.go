package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddress string
	DatabaseURL string
}

func Load() (Config, error) {
	port := envOrDefault("PORT", "8080")
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be an integer between 1 and 65535")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = buildDatabaseURL()
	}
	parsedDatabaseURL, err := url.ParseRequestURI(databaseURL)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DATABASE_URL: %w", err)
	}
	if parsedDatabaseURL.Scheme != "postgres" && parsedDatabaseURL.Scheme != "postgresql" {
		return Config{}, fmt.Errorf("DATABASE_URL must use postgres or postgresql scheme")
	}
	if parsedDatabaseURL.Host == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must contain a host")
	}

	return Config{
		HTTPAddress: ":" + port,
		DatabaseURL: databaseURL,
	}, nil
}

func buildDatabaseURL() string {
	databaseURL := &url.URL{
		Scheme: "postgres",
		Host:   envOrDefault("DB_HOST", "localhost") + ":" + envOrDefault("DB_PORT", "5432"),
		Path:   envOrDefault("DB_NAME", "persons"),
		User: url.UserPassword(
			envOrDefault("DB_USER", "program"),
			envOrDefault("DB_PASSWORD", "test"),
		),
	}
	query := databaseURL.Query()
	query.Set("sslmode", envOrDefault("DB_SSLMODE", "disable"))
	databaseURL.RawQuery = query.Encode()
	return databaseURL.String()
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
