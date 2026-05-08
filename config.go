package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	Port  string `mapstructure:"port"`
	DBURL string `mapstructure:"db_url"`
}

func LoadConfiguration() (config Configuration, err error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("Помилка читання .env: %w", err)
	}

	config.Port = viper.GetString("PORT")
	config.DBURL = viper.GetString("DB_URL")

	if config.Port == "" {
		config.Port = "8080"
	}

	return config, nil
}
