package lgtm

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GitHubGet(uri string) ([]byte, error) {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	token, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	url := fmt.Sprintf("https://api.github.com%s", uri)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}

	return body, nil
}

func GitHubAuthorize(user string, password string, fingerprint string, otpCode string) (*http.Response, error) {
	reqBody := []byte(fmt.Sprintf(`{"note":"lgtm","scopes":["repo"],"fingerprint":"%s"}`, fingerprint))

	req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)

	if otpCode != "" {
		req.Header.Add("X-GitHub-OTP", otpCode)
	}

	client := &http.Client{}
	return client.Do(req)
}
