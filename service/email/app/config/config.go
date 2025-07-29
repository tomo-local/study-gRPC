package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	// SMTP Server
	SMTP struct {
		Host         string
		Port         int
		TLS          bool
		AuthRequired bool
	}

	// IMAP Server
	IMAP struct {
		Host string
		Port int
		TLS  bool
	}

	// JWT
	JWT struct {
		Secret    string
		ExpiresIn time.Duration
	}

	// Application
	App struct {
		Env      string
		LogLevel string
	}
}

var AppConfig *Config

// Load configuration from environment variables
func Load() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig = &Config{}

	// Database configuration
	AppConfig.DB.Host = getEnv("DB_HOST", "localhost")
	AppConfig.DB.Port = getEnv("DB_PORT", "5432")
	AppConfig.DB.User = getEnv("DB_USER", "email_user")
	AppConfig.DB.Password = getEnv("DB_PASSWORD", "email_pass")
	AppConfig.DB.Name = getEnv("DB_NAME", "email_db")
	AppConfig.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	// SMTP configuration
	AppConfig.SMTP.Host = getEnv("SMTP_HOST", "localhost")
	AppConfig.SMTP.Port = getEnvAsInt("SMTP_PORT", 2525)
	AppConfig.SMTP.TLS = getEnvAsBool("SMTP_TLS", false)
	AppConfig.SMTP.AuthRequired = getEnvAsBool("SMTP_AUTH_REQUIRED", true)

	// IMAP configuration
	AppConfig.IMAP.Host = getEnv("IMAP_HOST", "localhost")
	AppConfig.IMAP.Port = getEnvAsInt("IMAP_PORT", 1143)
	AppConfig.IMAP.TLS = getEnvAsBool("IMAP_TLS", false)

	// JWT configuration
	AppConfig.JWT.Secret = getEnv("JWT_SECRET", "your-jwt-secret-key-here")
	AppConfig.JWT.ExpiresIn = getEnvAsDuration("JWT_EXPIRES_IN", 24*time.Hour)

	// Application configuration
	AppConfig.App.Env = getEnv("APP_ENV", "development")
	AppConfig.App.LogLevel = getEnv("LOG_LEVEL", "debug")

	log.Println("Configuration loaded successfully")
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultVal
}
