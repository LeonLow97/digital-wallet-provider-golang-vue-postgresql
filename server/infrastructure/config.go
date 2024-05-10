package infrastructure

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

var (
	MODE_TESTING       = "testing"
	MODE_DEVELOPMENT   = "development"
	MODE_DOCKER        = "docker"
	MODE_PREPRODUCTION = "preproduction"
	MODE_PRODUCTION    = "production"
)

type Environment struct {
	Mode string `mapstructure:"mode"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Issuer string `mapstructure:"issuer"`
}

type SMTPConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type Config struct {
	Env  Environment `mapstructure:"environment"`
	JWT  JWTConfig   `mapstructure:"jwt"`
	SMTP SMTPConfig  `mapstructure:"smtp"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("mode", MODE_DEVELOPMENT)

	viper.SetConfigName("development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("config file not found", err)
			return nil, errors.New("config file not found")
		} else {
			// Config file was found but another error was produced
			log.Println("error reading config file", err)
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
