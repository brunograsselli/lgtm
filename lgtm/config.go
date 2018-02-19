package lgtm

import (
	"github.com/spf13/viper"
)

type Config struct {
	UserName string
	Repos    []string
}

func NewConfig() *Config {
	return &Config{
		UserName: viper.GetString("username"),
		Repos:    viper.GetStringSlice("repos"),
	}
}
