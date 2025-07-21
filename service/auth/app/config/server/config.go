package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Mailer   MailerConfig   `yaml:"mailer"`
}

type ServerConfig struct {
	Port int `yaml:"port" envconfig:"SERVER_PORT"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" envconfig:"DB_HOST" default:"localhost"`
	Port     int    `yaml:"port" envconfig:"DB_PORT" default:"5432"`
	User     string `yaml:"user" envconfig:"DB_USER" default:"postgres"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD" default:"password"`
	Name     string `yaml:"name" envconfig:"DB_NAME" default:"auth_db"`
	SSLMode  string `yaml:"ssl_mode" envconfig:"DB_SSL_MODE" default:"disable"`
}

type JWTConfig struct {
	SecretKey      string `yaml:"secret_key" envconfig:"JWT_SECRET_KEY"`
	TokenDurationH int    `yaml:"token_duration_hours" envconfig:"JWT_TOKEN_DURATION_HOURS"`
}

type MailerConfig struct {
	Host     string `yaml:"host" envconfig:"MAILER_HOST"`
	Port     int    `yaml:"port" envconfig:"MAILER_PORT"`
	Username string `yaml:"username" envconfig:"MAILER_USERNAME"`
	Password string `yaml:"password" envconfig:"MAILER_PASSWORD"`
	From     string `yaml:"from" envconfig:"MAILER_FROM"`
}

// LoadConfigWithAutoPath は設定ファイルのパスを自動決定して設定を読み込む
func LoadConfigWithAutoPath() (*Config, error) {
	configPath := findConfigPath()
	log.Printf("Loading config from: %s", configPath)
	return LoadConfig(configPath)
}

// findConfigPath は設定ファイルのパスを決定する
func findConfigPath() string {
	// 1. 環境変数をチェック
	if configPath := os.Getenv("AUTH_CONFIG_PATH"); configPath != "" {
		log.Printf("Using config path from environment variable: %s", configPath)
		return configPath
	}

	// 2. 実行ファイルの場所から相対的にパスを計算
	execPath, err := os.Executable()
	if err == nil {
		// 実行ファイルのディレクトリから ../../config.yml を計算
		execDir := filepath.Dir(execPath)
		configPath := filepath.Join(execDir, "..", "..", "config.yml")
		if _, err := os.Stat(configPath); err == nil {
			log.Printf("Found config file relative to executable: %s", configPath)
			return configPath
		}
	}

	// 3. デフォルトの相対パス
	defaultPaths := []string{
		"../../config.yml", // 開発時の場合
		"../config.yml",    // 別の場合
		"./config.yml",     // 同じディレクトリ
		"config.yml",       // カレントディレクトリ
	}

	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Found config file at: %s", path)
			return path
		}
	}

	// 4. デフォルト値を返す（警告付き）
	defaultPath := "../../config.yml"
	log.Printf("Warning: Config file not found, using default path: %s", defaultPath)
	return defaultPath
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

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

	// envconfigを使って環境変数で設定を上書き
	if err := envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	return config, nil
}
