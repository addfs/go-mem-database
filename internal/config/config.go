package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Logger struct {
		Level  zap.AtomicLevel `yaml:"level" mapstructure:"level"`
		Output string          `yaml:"output" mapstructure:"output"`
	} `yaml:"logger" mapstructure:"logger"`
	EngineConfig struct {
		Type string `yaml:"type" mapstructure:"type"`
	} `yaml:"engine" mapstructure:"engine"`
	Network struct {
		Address        string `yaml:"address" mapstructure:"address"`
		MaxConnections int    `yaml:"max_connections" mapstructure:"max_connections"`
	}
}

func NewConfig(configFile string) func() *Config {
	return func() *Config {
		v := viper.New()

		v.SetConfigFile(configFile)

		v.SetDefault("logger.level", "info")

		v.SetDefault("engine.type", "in_memory")

		v.SetDefault("network.address", "127.0.0.1:3223")
		v.SetDefault("network.max_connections", 100)

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
