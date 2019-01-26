package lgtm

import (
	"io/ioutil"
	"os"
)

type Secrets struct {
	path string
}

func NewSecrets(path string) *Secrets {
	return &Secrets{path: path}
}

func (s *Secrets) Token() ([]byte, error) {
	return ioutil.ReadFile(s.path)
}

func (s *Secrets) CheckToken() bool {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *Secrets) SaveToken(token []byte) error {
	return ioutil.WriteFile(s.path, token, 0644)
}

func (s *Secrets) DeleteToken() error {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(s.path)
}
