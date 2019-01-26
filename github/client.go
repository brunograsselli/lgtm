package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token []byte
}

func NewClient(token []byte) *Client {
	return &Client{token: token}
}

func Authorize(user string, password string, fingerprint string, otpCode string) (*http.Response, error) {
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

func (c *Client) PullRequests(repository string) ([]PullRequest, error) {
	body, err := c.get(fmt.Sprintf("/repos/%s/pulls?sort=created&direction=asc", repository))

	if err != nil {
		return nil, err
	}

	var pullRequests []PullRequest
	json.Unmarshal([]byte(body), &pullRequests)

	return pullRequests, nil
}

func (c *Client) Reviews(repository string, pullRequestNumber int32) ([]Review, error) {
	body, err := c.get(fmt.Sprintf("/repos/%s/pulls/%d/reviews?sort=created&direction=asc", repository, pullRequestNumber))

	if err != nil {
		return nil, err
	}

	var reviews []Review
	json.Unmarshal([]byte(body), &reviews)

	return reviews, nil
}

func (c *Client) get(uri string) ([]byte, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.github.com%s", uri)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", c.token))

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
