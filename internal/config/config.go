package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Logger struct {
		Level zap.AtomicLevel `yaml:"level" mapstructure:"level"`
	} `yaml:"logger" mapstructure:"logger"`
	EngineConfig struct {
		Type string `yaml:"type" mapstructure:"type"`
	} `yaml:"engine" mapstructure:"engine"`
}

func NewConfig(configFile string) func() *Config {
	return func() *Config {
		v := viper.New()

		v.SetConfigFile(configFile)

		v.SetDefault("logger.level", "info")
		v.SetDefault("engine.type", "in_memory")

		if err := v.ReadInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				panic(err)
			}
		}

		config := new(Config)

		decodeHook := mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.TextUnmarshallerHookFunc(),
		)

		if err := v.Unmarshal(config, viper.DecodeHook(decodeHook)); err != nil {
			panic(err)
		}

		return config
	}
}
