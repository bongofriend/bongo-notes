package config

import (
	"os"

	"github.com/spf13/viper"
)

type Environment string

const (
	local Environment = "local"
)

type Config struct {
	Port        int
	localConfig localConfig
}

type localConfig struct {
	Port           int
	IncludeSwagger bool
}

func (c Config) IncludeSwagger() bool {
	return c.localConfig.IncludeSwagger
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
	var c localConfig
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &Config{
		Port:        c.Port,
		localConfig: c,
	}, nil
}
