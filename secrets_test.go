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

	s := &Secrets{Path: path}

	token, err := s.Token()

	if err != nil {
		t.Error(err)
	}

	if string(token) != "fake_token" {
		t.Errorf("Got '%s', want 'fake_token'", token)
	}
}
