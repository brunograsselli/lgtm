package lgtm

import (
	"github.com/spf13/viper"
)

type Config struct {
}

func (c *Config) Repos() []string {
	return viper.GetStringSlice("repos")
}

func (c *Config) UserName() string {
	return viper.GetString("username")
}
