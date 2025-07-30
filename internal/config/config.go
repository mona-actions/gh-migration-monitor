package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	GitHub struct {
		Token        string `mapstructure:"token"`
		Organization string `mapstructure:"organization"`
	} `mapstructure:"github"`

	Migration struct {
		IsLegacy bool `mapstructure:"is_legacy"`
	} `mapstructure:"migration"`

	Output struct {
		Format string `mapstructure:"format"`
		Quiet  bool   `mapstructure:"quiet"`
	} `mapstructure:"output"`
}

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gh-migration-monitor")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("GHMM")
	viper.AutomaticEnv()

	// Bind specific environment variables
	viper.BindEnv("github.token", "GHMM_GITHUB_TOKEN")
	viper.BindEnv("github.organization", "GHMM_GITHUB_ORGANIZATION")
	viper.BindEnv("migration.is_legacy", "GHMM_ISLEGACY")

	// Read configuration file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.GitHub.Organization == "" {
		return fmt.Errorf("github organization is required")
	}

	if c.GitHub.Token == "" {
		return fmt.Errorf("github token is required")
	}

	return nil
}
