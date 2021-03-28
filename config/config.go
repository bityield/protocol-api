package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	// Name for the applicaiton
	Name = "protocol-api"

	// Env getter for Viper
	Env = "env"

	// Environment getter for Viper
	Environment = "environment"

	// Development environment
	Development = "development"

	// Dev environment
	Dev = "dev"

	// Stage environment
	Stage = "stage"

	// Prod environment
	Prod = "prod"

	// Port variable
	Port = "port"

	environmentVar     = "ENVIRONMENT"
	environmentDefault = Development

	portVar     = "PORT"
	portDefault = "8000"
)

// Reader represents configuration reader
type Reader interface {
	Get(string) interface{}
	GetString(string) string
	GetInt(string) int
	GetBool(string) bool
	GetStringMap(string) map[string]interface{}
	GetStringMapString(string) map[string]string
	GetStringSlice(string) []string
	SetDefault(string, interface{})
}

// DefaultSettings is the function for configuring defaults
type DefaultSettings func(config Reader)

// MainDefaults - returns the defauls of the config passed
// func MainDefaults(config Reader) {
// 	Defaults(config)
// }

// Defaults is the default settings functor
func Defaults(config Reader) {
	config.SetDefault(Environment, GetEnv(environmentVar, Development))
	config.SetDefault(Port, GetEnv(portVar, portDefault))
}

// GetEnv - pull values or set defaults.
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}

// LoadConfig - returns configuration for a particular app
func LoadConfig(defaultSetup DefaultSettings) (Reader, error) {
	config := viper.New()

	// Configure the defaults
	Defaults(config)

	// if config.GetString(Environment) == Development || config.GetString(Environment) == Prod {
	// 	panic("missing or invalid [ENVIRONMENT] setting")
	// }

	return config, nil
}

// LoadLogger - set the defaults for the logging class
func LoadLogger(config Reader) *logrus.Logger {
	log := logrus.New()

	log.Formatter = new(prefixed.TextFormatter)

	log.Out = os.Stdout

	// log.SetLevel(logrus.WarnLevel)
	log.SetLevel(logrus.InfoLevel)

	return log
}
