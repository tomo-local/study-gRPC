package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Port     int      `mapstructure:"port" default:"8082"`
	Env      string   `mapstructure:"env" default:"development"`
	Database Database `mapstructure:"database"`
}

// Database holds the database configuration
type Database struct {
	Host     string `mapstructure:"host" default:"localhost"`
	Port     int    `mapstructure:"port" default:"5432"`
	User     string `mapstructure:"user" default:"noteuser"`
	Password string `mapstructure:"password" default:"notepass"`
	DBName   string `mapstructure:"dbname" default:"notedb"`
	SSLMode  string `mapstructure:"sslmode" default:"disable"`
}

// LoadConfig loads the configuration from a file or environment variables
func LoadConfig() (*Config, error) {
	// Set the configuration file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set the path to look for the configuration file
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Join(filepath.Dir(filename), "../../../")
	viper.AddConfigPath(configPath)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the configuration into a Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Override port from settings if provided
	if port := viper.GetInt("settings.PORT"); port != 0 {
		config.Port = port
	}

	// Override env from settings if provided
	if env := viper.GetString("settings.ENV"); env != "" {
		config.Env = env
	}

	return &config, nil
}
