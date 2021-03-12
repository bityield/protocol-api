package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Config Reader

func init() {
	os.Setenv("ENVIRONMENT", "test")
	Config, _ = LoadConfig(Defaults)
	os.Setenv("ENVIRONMENT", "")

	if Config.GetString("environment") != "test" {
		panic(fmt.Errorf("test [environment] is not in [test] mode"))
	}
}

func TestConstants(t *testing.T) {
	assert.Equal(t, Env, "env")
	assert.Equal(t, Environment, "environment")
	assert.Equal(t, environmentVar, "ENVIRONMENT")
	assert.Equal(t, environmentDefault, "development")

	assert.Equal(t, Dev, "dev")
	assert.Equal(t, Stage, "stage")
	assert.Equal(t, Prod, "prod")

	assert.Equal(t, Port, "port")
	assert.Equal(t, portVar, "PORT")
	assert.Equal(t, portDefault, "8000")
}

func TestLoadLogger(t *testing.T) {
	logger := LoadLogger(Config)

	if assert.NotNil(t, logger) {
		assert.Equal(t, logger.Formatter, &prefixed.TextFormatter{})
		assert.Equal(t, logger.Level, logrus.WarnLevel)
	}
}

func TestGetEnvExists(t *testing.T) {
	os.Setenv("FOO", "nothing")

	assert.Equal(t, GetEnv("FOO", "invalid"), "nothing")

	os.Unsetenv("FOO")
}

func TestGetEnvNotExists(t *testing.T) {
	os.Setenv("FOO", "")

	assert.Equal(t, GetEnv("FOO", "invalid"), "invalid")

	os.Unsetenv("FOO")
}
