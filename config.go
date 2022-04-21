package main

import (
	lr "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type arahServerConfig struct {
	Server struct {
		Port   int    `mapstructure:"port"`
		Domain string `mapstructure:"domain"`
	} `mapstructure:"server"`
	Captcha struct {
		Sitekey string `mapstructure:"sitekey"`
		Secret  string `mapstructure:"secret"`
	}
	Mongo struct {
		URI      string `mapstructure:"uri"`
		Database string `mapstructure:"database"`
	}
}

func loadConfig() (res *arahServerConfig) {
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
