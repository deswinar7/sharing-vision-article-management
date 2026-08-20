package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv             string
	Port               string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	DBTLSMode          string
	DBCACertPath       string
	CORSAllowedOrigins []string
	ShutdownTimeout    time.Duration
}

func Load() (Config, error) {
	timeoutSeconds, err := strconv.Atoi(env("SHUTDOWN_TIMEOUT_SECONDS", "10"))
	if err != nil || timeoutSeconds < 1 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT_SECONDS must be a positive integer")
	}

	return Config{
		AppEnv:             env("APP_ENV", "development"),
		Port:               env("PORT", "8080"),
		DBHost:             env("DB_HOST", "localhost"),
		DBPort:             env("DB_PORT", "3306"),
		DBName:             env("DB_NAME", "article"),
		DBUser:             env("DB_USER", "article_user"),
		DBPassword:         env("DB_PASSWORD", "article_password"),
		DBTLSMode:          env("DB_TLS_MODE", "false"),
		DBCACertPath:       os.Getenv("DB_CA_CERT_PATH"),
		CORSAllowedOrigins: splitCSV(env("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		ShutdownTimeout:    time.Duration(timeoutSeconds) * time.Second,
	}, nil
}

func (c Config) Address() string { return "0.0.0.0:" + c.Port }

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
