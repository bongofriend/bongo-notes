package config

import (
	"os"

	"github.com/spf13/viper"
)

type Environment string

type Config struct {
	Port int `mapstruture:"port"`
	Db   struct {
		Driver string `mapstructure:"driver"`
		Path   string `mapstructure:"path"`
	} `mapstructure:"db"`
	IncludeSwagger  bool
	JwtSecret       string `mapstruture:"jwtSecret"`
	NotesFolderPath string `mapstructure:"notesPath"`
}

func LoadConfig(configPath string) (Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return Config{}, err
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
