package config

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	LogLevel string `mapstructure:"logLevel"`
}

type BrokerConfig struct {
	Uri string `mapstructure:"uri"`
}

// LoadConfig takes a config struct of the calling microservice.
func LoadConfig(configStruct interface{}) {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("csgo")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	if err := viper.Unmarshal(&configStruct); err != nil {
		log.Panic(err)
	}
}

// IsDebug returns whether the application is in debug mode.
func (c *GlobalConfig) IsDebug() bool {
	return c.LogLevel == "debug" || c.IsTrace()
}

// IsTrace returns whether the application should do extended debugging.
func (c *GlobalConfig) IsTrace() bool {
	return c.LogLevel == "trace"
}

// SetLoggingLevel sets the logging level in relation to the level set in the config file.
func (c *GlobalConfig) SetLoggingLevel() {
	if c.IsTrace() {
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	} else if c.IsDebug() {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
