package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Mailer   MailerConfig   `yaml:"mailer"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type JWTConfig struct {
	SecretKey      string `yaml:"secret_key"`
	TokenDurationH int    `yaml:"token_duration_hours"`
}

type MailerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// デフォルト値の設定
	setDefaults(config)

	// 設定ファイルを読み込み
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			file, err := os.Open(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open config file: %w", err)
			}
			defer file.Close()

			decoder := yaml.NewDecoder(file)
			if err := decoder.Decode(config); err != nil {
				return nil, fmt.Errorf("failed to decode config file: %w", err)
			}
		}
	}

	// 環境変数で設定を上書き
	if err := overrideWithEnv(config); err != nil {
		return nil, fmt.Errorf("failed to override config with environment variables: %w", err)
	}

	return config, nil
}

func overrideWithEnv(config *Config) error {
	// Server
	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		config.Server.Port = p
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
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// JWT
	if secretKey := os.Getenv("JWT_SECRET_KEY"); secretKey != "" {
		config.JWT.SecretKey = secretKey
	}
	if tokenDuration := os.Getenv("JWT_TOKEN_DURATION_HOURS"); tokenDuration != "" {
		d, err := strconv.Atoi(tokenDuration)
		if err != nil {
			return fmt.Errorf("invalid JWT_TOKEN_DURATION_HOURS: %w", err)
		}
		config.JWT.TokenDurationH = d
	}

	// Mailer
	if host := os.Getenv("MAILER_HOST"); host != "" {
		config.Mailer.Host = host
	}
	if port := os.Getenv("MAILER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid MAILER_PORT: %w", err)
		}
		config.Mailer.Port = p
	}
	if username := os.Getenv("MAILER_USERNAME"); username != "" {
		config.Mailer.Username = username
	}
	if password := os.Getenv("MAILER_PASSWORD"); password != "" {
		config.Mailer.Password = password
	}
	if from := os.Getenv("MAILER_FROM"); from != "" {
		config.Mailer.From = from
	}

	return nil
}

func setDefaults(config *Config) {
	// Server
	if config.Server.Port == 0 {
		config.Server.Port = 50053
	}

	// Database
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.User == "" {
		config.Database.User = "postgres"
	}
	if config.Database.Name == "" {
		config.Database.Name = "auth_db"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// JWT
	if config.JWT.SecretKey == "" {
		config.JWT.SecretKey = "your-secret-key-change-this-in-production"
	}
	if config.JWT.TokenDurationH == 0 {
		config.JWT.TokenDurationH = 24
	}

	// Mailer
	if config.Mailer.Host == "" {
		config.Mailer.Host = "smtp.gmail.com"
	}
	if config.Mailer.Port == 0 {
		config.Mailer.Port = 587
	}
}
