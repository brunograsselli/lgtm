package lgtm

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

var configPath = fmt.Sprintf("%s/.lgtm.yml", os.Getenv("HOME"))

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

func (c *Config) SaveUserName(username string) error {
	c.UserName = username

	y, err := yaml.Marshal(c)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, y, 0644)
}
