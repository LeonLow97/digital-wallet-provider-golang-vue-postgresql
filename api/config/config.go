package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Development Development `mapstructure:"development"`
	Production  Production  `mapstructure:"production"`
}

type Development struct {
	URL string `mapstructure:"url"`
}

type Production struct {
	URL string `mapstructure:"url"`
}


func LoadConfig(environment string) (*Config, error) {
	// using `viper` library to simplify configuration management
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// viper.AddConfigPath(".") // will look for the config file in the current directory
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
