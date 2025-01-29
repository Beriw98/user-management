package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURI string
}

func New() *Config {

	vpr := viper.New()
	vpr.AutomaticEnv()
	vpr.SetConfigName("config")
	vpr.AddConfigPath(".")
	vpr.SetConfigType("yaml")

	if err := vpr.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			panic("error reading config file")
		}
	}

	return &Config{
		DatabaseURI: vpr.GetString("database_uri"),
	}
}
