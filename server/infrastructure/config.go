package infrastructure

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	MODE_TESTING       = "testing"
	MODE_DEVELOPMENT   = "development"
	MODE_DOCKER        = "docker"
	MODE_PREPRODUCTION = "preproduction"
	MODE_PRODUCTION    = "production"
)

type Config struct {
	Mode   string `mapstructure:"mode"`
	Server struct {
		Port       int    `mapstructure:"port"`
		BackendURL string `mapstructure:"backend_url"`
	} `mapstructure:"server"`
	Frontend struct {
		FrontendURL string `mapstructure:"frontend_url"`
	} `mapstructure:"frontend"`
	JWT struct {
		Secret string `mapstructure:"secret"`
		Issuer string `mapstructure:"issuer"`
	} `mapstructure:"jwt"`
	SMTP struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"smtp"`
	TOTP struct {
		Issuer        string `mapstructure:"issuer"`
		EncryptionKey string `mapstructure:"encryption_key"`
	} `mapstructure:"totp"`
	CSRF struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"csrf"`
	Postgres struct {
		PostgresUser     string `mapstructure:"postgres_user"`
		PostgresPassword string `mapstructure:"postgres_password"`
		PostgresHost     string `mapstructure:"postgres_host"`
		PostgresPort     int    `mapstructure:"postgres_port"`
		PostgresDB       string `mapstructure:"postgres_db"`
	} `mapstructure:"postgres"`
	Redis struct {
		RedisHost string `mapstructure:"redis_host"`
		RedisPort int    `mapstructure:"redis_port"`
	} `mapstructure:"redis"`
}

func LoadConfig() (*Config, error) {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development"
	}

	viper.SetConfigName(mode)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// inject environment variables to override matching keys in configuration files (.yaml)
	viper.AutomaticEnv()

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
