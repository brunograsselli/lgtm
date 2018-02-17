package lgtm

import (
	"github.com/spf13/viper"
)

type Repos struct {
}

func (r *Repos) All() []string {
	return viper.GetStringSlice("repos")
}
