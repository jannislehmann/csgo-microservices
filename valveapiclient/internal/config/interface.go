package config

// ConfigUseCase defines the config service functions.
type ConfigUseCase interface {
	GetConfig() *Config
}
