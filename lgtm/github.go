package lgtm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GitHubGet(uri string, secrets *Secrets) ([]byte, error) {
	token, err := secrets.Token()

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

func GitHubPullRequests(repository string, secrets *Secrets) ([]PullRequest, error) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls?sort=created&direction=asc", repository), secrets)

	if err != nil {
		return nil, err
	}

	var pullRequests []PullRequest
	json.Unmarshal([]byte(body), &pullRequests)

	return pullRequests, nil
}
