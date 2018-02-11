package lgtm

import (
	"io/ioutil"
	"os"
)

type Secrets struct {
	Path string
}

func (s *Secrets) Token() ([]byte, error) {
	return ioutil.ReadFile(s.Path)
}

func (s *Secrets) CheckToken() bool {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *Secrets) SaveToken(token string) error {
	return ioutil.WriteFile(s.Path, []byte(token), 0644)
}

func (s *Secrets) DeleteToken() error {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(s.Path)
}
