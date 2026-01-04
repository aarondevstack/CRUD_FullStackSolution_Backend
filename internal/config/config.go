package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	API      APIConfig      `mapstructure:"api"`
	Database DatabaseConfig `mapstructure:"database"`
}

type APIConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type       string           `mapstructure:"type"`
	Host       string           `mapstructure:"host"`
	Port       int              `mapstructure:"port"`
	User       string           `mapstructure:"user"`
	Password   string           `mapstructure:"password"`
	Name       string           `mapstructure:"name"`
	AutoBackup AutoBackupConfig `mapstructure:"autoBackup"`
}

type AutoBackupConfig struct {
	Enable bool   `mapstructure:"enable"`
	Cron   string `mapstructure:"cron"`
}

var AppConfig *Config

// Load loads configuration from file or embedded config
func Load() error {
	v := viper.New()
	v.SetConfigType("yaml")

	// Try to load from config/config.yaml in current directory
	configPath := filepath.Join("config", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Load from embedded config
		if err := v.ReadConfig(bytes.NewReader(embeddedConfig)); err != nil {
			return fmt.Errorf("failed to read embedded config: %w", err)
		}
	}

	AppConfig = &Config{}
	if err := v.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// GetBaseDir returns the directory where the executable is located
// This is where MySQL and other runtime files will be stored
func GetBaseDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	// Return the directory containing the executable
	return filepath.Dir(execPath), nil
}
