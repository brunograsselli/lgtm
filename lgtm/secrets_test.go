package lgtm

import (
	"io/ioutil"
	"os"
	"testing"
)

var path = "/tmp/lgtm-test.secret"

func TestTokenWhenFileIsPresent(t *testing.T) {
	ioutil.WriteFile(path, []byte("fake_token"), 0644)

	defer os.Remove(path)

	s := NewSecrets(path)

	token, err := s.Token()

	if err != nil {
		t.Error(err)
	}

	if string(token) != "fake_token" {
		t.Errorf("Got '%s', want 'fake_token'", token)
	}
}

func TestCheckWhenFileIsPresent(t *testing.T) {
	ioutil.WriteFile(path, []byte("fake_token"), 0644)

	defer os.Remove(path)

	s := NewSecrets(path)

	result := s.CheckToken()

	if result != true {
		t.Errorf("Got %v, want true", result)
	}
}

func TestCheckWhenFileIsNotPresent(t *testing.T) {
	s := NewSecrets(path)

	result := s.CheckToken()

	if result != false {
		t.Errorf("Got %v, want false", result)
	}
}

func TestSaveToken(t *testing.T) {
	s := NewSecrets(path)

	err := s.SaveToken([]byte("abc"))

	if err != nil {
		t.Error(err)
	}

	token, err := ioutil.ReadFile(path)

	if err != nil {
		t.Error(err)
	}

	if string(token) != "abc" {
		t.Errorf("Got '%s', want 'abc'", token)
	}
}

func TestDeleteTokenWhenFileIsPresent(t *testing.T) {
	ioutil.WriteFile(path, []byte("fake_token"), 0644)

	defer func() {
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}()

	s := NewSecrets(path)

	err := s.DeleteToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(path); err == nil {
		t.Error("Expect secret file to not exist but it exists")
	}
}

func TestDeleteTokenWhenFileIsNotPresent(t *testing.T) {
	s := NewSecrets(path)

	err := s.DeleteToken()

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(path); err == nil {
		t.Error("Expect secret file to not exist but it exists")
	}
}
