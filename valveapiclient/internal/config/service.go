package config

import (
	shared_config "github.com/Cludch/csgo-microservices/shared/pkg/config"
)

type ConfigService struct {
	config *Config
}

func NewService() *ConfigService {
	service := ConfigService{
		config: &Config{},
	}
	shared_config.LoadConfig(&service.config)

	service.config.Global.SetLoggingLevel()

	return &service
}

// Config holds the application configuration.
type Config struct {
	Global   *shared_config.GlobalConfig `mapstructure:"global"`
	Broker   *shared_config.BrokerConfig `mapstructure:"broker"`
	Database *DatabaseConfig             `mapstructure:"database"`
}

// DatabaseConfig holds database connection information.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// GetConfig returns the application configuration.
func (s *ConfigService) GetConfig() *Config {
	return s.config
}
