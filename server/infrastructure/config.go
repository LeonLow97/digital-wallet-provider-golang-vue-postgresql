package bootstrap

import "github.com/spf13/viper"

type Config struct {
	
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
