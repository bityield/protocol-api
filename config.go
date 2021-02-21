package main

import (
	"os"

	"github.com/spf13/viper"
)

const (
	// Name for the applicaiton
	Name = "bityield-api"

	// Env getter for Viper
	Env = "env"

	// Environment getter for Viper
	Environment = "environment"

	// EnvironmentVar export env variable
	EnvironmentVar = "ENVIRONMENT"

	// Local environment
	Local = "local"

	// Production environment
	Production = "production"
)

// ConfigReader represents configuration reader
type ConfigReader interface {
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
type DefaultSettings func(config ConfigReader)

// ConfigDefaults - returns the defauls of the config passed
func ConfigDefaults(config ConfigReader) {
	Defaults(config)
}

// Defaults is the default settings functor
func Defaults(config ConfigReader) {
	config.SetDefault(Environment, GetEnv(EnvironmentVar, Local))
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
func LoadConfig(defaultSetup DefaultSettings) (ConfigReader, error) {
	config := viper.New()

	// Configure the defaults
	Defaults(config)

	if config.GetString(Environment) == Local || config.GetString(Environment) == Production {
		panic("missing or invalid [ENVIRONMENT] setting")
	}

	return config, nil
}
