package main

import (
	lr "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ArahServerConfig struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func loadConfig() (res *ArahServerConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		lr.Fatalf("error reading config : %s", err)
	}

	err = viper.Unmarshal(&res)

	if err != nil {
		lr.Fatalf("error unmarshalling config : %s", err)
	}

	return
}
