package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type LoggerConfig struct {
	Level  string `envconfig:"LOGGER_LEVEL" required:"true"`
	Folder string `envconfig:"LOGGER_FOLDER" required:"true"`
}

func NewConfig() (LoggerConfig, error) {
	var config LoggerConfig

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return LoggerConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() LoggerConfig {
	var config LoggerConfig

	if err := envconfig.Process("", &config); err != nil {
		err = fmt.Errorf("process envconfig: %w", err)
		panic(err)
	}

	return config
}