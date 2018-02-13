package lgtm

import (
	"fmt"
	"sort"

	"github.com/spf13/viper"
)

type Repos struct {
}

func (r *Repos) List() {
	all := r.All()

	sort.Strings(all)

	for _, repo := range all {
		fmt.Println(repo)
	}
}

func (r *Repos) All() []string {
	return viper.GetStringSlice("repos")
}
