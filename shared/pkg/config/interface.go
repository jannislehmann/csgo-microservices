package config

// UseCase defines the config service functions.
type UseCase interface {
	LoadConfig(interface{}) interface{}
}
