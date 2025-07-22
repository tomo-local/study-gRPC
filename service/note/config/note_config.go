package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     int               `yaml:"port"`
	Env      string            `yaml:"env"`
	Database DatabaseConfig    `yaml:"database"`
	Auth     AuthServiceConfig `yaml:"auth"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// AuthServiceConfig はAuth Serviceとの連携設定
type AuthServiceConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	URL            string `yaml:"url"`
	Enabled        bool   `yaml:"enabled"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	JWTSecretKey   string `yaml:"jwt_secret_key"`
}

// GetAuthServiceAddress は完全なAuth ServiceのアドレスURL取得
func (c *AuthServiceConfig) GetAuthServiceAddress() string {
	if c.URL != "" {
		return c.URL
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetTimeout はAuth Serviceへのリクエストタイムアウトを取得
func (c *AuthServiceConfig) GetTimeout() time.Duration {
	if c.TimeoutSeconds <= 0 {
		return 10 * time.Second // デフォルト10秒
	}
	return time.Duration(c.TimeoutSeconds) * time.Second
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	config := &Config{}

	// デフォルト値を設定
	setDefaults(config)

	// 環境変数で設定を上書き
	if err := overrideWithEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults(config *Config) {
	// Server
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.Env == "" {
		config.Env = "development"
	}

	// Database defaults
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.User == "" {
		config.Database.User = "noteuser"
	}
	if config.Database.DBName == "" {
		config.Database.DBName = "notedb"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// Auth Service defaults
	if config.Auth.Host == "" {
		config.Auth.Host = "localhost"
	}
	if config.Auth.Port == 0 {
		config.Auth.Port = 9001
	}
	if config.Auth.TimeoutSeconds == 0 {
		config.Auth.TimeoutSeconds = 10
	}
	config.Auth.Enabled = true
}

// overrideWithEnv overrides configuration with environment variables
func overrideWithEnv(config *Config) error {
	// Server
	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		config.Port = p
	}
	if env := os.Getenv("ENV"); env != "" {
		config.Env = env
	}

	// Database
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid DB_PORT: %w", err)
		}
		config.Database.Port = p
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSL_MODE"); sslmode != "" {
		config.Database.SSLMode = sslmode
	}

	// Auth Service
	if host := os.Getenv("AUTH_SERVICE_HOST"); host != "" {
		config.Auth.Host = host
	}
	if port := os.Getenv("AUTH_SERVICE_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid AUTH_SERVICE_PORT: %w", err)
		}
		config.Auth.Port = p
	}
	if url := os.Getenv("AUTH_SERVICE_URL"); url != "" {
		config.Auth.URL = url
	}
	if enabled := os.Getenv("AUTH_ENABLED"); enabled != "" {
		e, err := strconv.ParseBool(enabled)
		if err != nil {
			return fmt.Errorf("invalid AUTH_ENABLED: %w", err)
		}
		config.Auth.Enabled = e
	}
	if timeout := os.Getenv("AUTH_TIMEOUT_SECONDS"); timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			return fmt.Errorf("invalid AUTH_TIMEOUT_SECONDS: %w", err)
		}
		config.Auth.TimeoutSeconds = t
	}
	if jwtSecret := os.Getenv("JWT_SECRET_KEY"); jwtSecret != "" {
		config.Auth.JWTSecretKey = jwtSecret
	}

	return nil
}
