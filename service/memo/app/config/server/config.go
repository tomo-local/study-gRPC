package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Port       int    `mapstructure:"port" default:"8080"`
	Env        string `mapstructure:"env" default:"development"`
	FolderPath string `mapstructure:"folder_path" default:"/tmp/memo"`
}

// EnvVar は env 配列の1要素を表す構造体だよ！
type EnvVar struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
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

	if port := viper.GetInt("settings.PORT"); port != 0 {
		config.Port = port
	}

	if env := viper.GetString("settings.ENV"); env != "" {
		config.Env = env
	}

	if folderPath := viper.GetString("settings.FOLDER_PATH"); folderPath != "" {
		config.FolderPath = folderPath
	}

	return &config, nil
}
