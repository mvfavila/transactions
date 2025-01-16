package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogFile  string `yaml:"log_file"`
	Port     string `yaml:"port"`
	Database struct {
		Driver string `yaml:"driver"`
		Source string `yaml:"source"`
	} `yaml:"database"`
	ExpectedDateFormat string `yaml:"expected_date_format"`
	TreasuryAPIBaseURL string `yaml:"treasury_api_base_url"`
}

var (
	AppConfig *Config
	once      sync.Once
)

// LoadConfig loads the configuration for the given environment.
// If AppConfig is already set, it skips reloading.
func LoadConfig(env string) error {
	var err error
	once.Do(func() {
		AppConfig, err = loadConfigFromFile(env)
	})
	return err
}

// loadConfigFromFile reads the config file and returns the configuration.
func loadConfigFromFile(env string) (*Config, error) {
	fileName := fmt.Sprintf("config/%s.yaml", env)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}

// LoadDefaultConfig initializes AppConfig with a default config for testing purposes.
func LoadDefaultConfig() {
	AppConfig = &Config{
		LogFile: "test.log",
		Port:    "8080",
		Database: struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		}{
			Driver: "sqlite3",
			Source: ":memory:",
		},
		ExpectedDateFormat: "2006-01-02",
		TreasuryAPIBaseURL: "https://api.fiscaldata.treasury.gov/services/api/test",
	}
}
