package config

import (
	"os"

	"github.com/spf13/viper"
)

type Environment string

type Config struct {
	Port int
	Db   struct {
		Driver string
		Path   string
	}
	IncludeSwagger bool
}

func LoadConfig(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
