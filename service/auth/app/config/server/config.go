package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Mailer   MailerConfig   `yaml:"mailer"`
}

type ServerConfig struct {
	Port int `yaml:"port" envconfig:"SERVER_PORT" default:"8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" envconfig:"DB_HOST" required:"true"`
	Port     int    `yaml:"port" envconfig:"DB_PORT" default:"5432"`
	User     string `yaml:"user" envconfig:"DB_USER" required:"true"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD" required:"true"`
	Name     string `yaml:"name" envconfig:"DB_NAME" required:"true"`
	SSLMode  string `yaml:"ssl_mode" envconfig:"DB_SSL_MODE" default:"disable"`
}

type JWTConfig struct {
	SecretKey      string `yaml:"secret_key" envconfig:"JWT_SECRET_KEY" required:"true"`
	TokenDurationH int    `yaml:"token_duration_hours" envconfig:"JWT_TOKEN_DURATION_HOURS" default:"24"`
}

type MailerConfig struct {
	Host     string `yaml:"host" envconfig:"MAILER_HOST" required:"true"`
	Port     int    `yaml:"port" envconfig:"MAILER_PORT" required:"true"`
	Username string `yaml:"username" envconfig:"MAILER_USERNAME" required:"true"`
	Password string `yaml:"password" envconfig:"MAILER_PASSWORD" required:"true"`
	From     string `yaml:"from" envconfig:"MAILER_FROM" required:"true"`
}

// LoadConfig 環境変数のみから設定を読み込む
func LoadConfig() (*Config, error) {
	config := &Config{}

	// 環境変数から設定を読み込み
	if err := envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	return config, nil
}
