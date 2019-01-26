package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	protocol        = "https://"
	host            = "api.github.com"
	authURI         = "/authorizations"
	pullRequestsURI = "/repos/%s/pulls?sort=created&direction=asc"
	reviewsURI      = "/repos/%s/pulls/%d/reviews?sort=created&direction=asc"

	httpStatusOk           = 200
	httpStatusCreated      = 201
	httpStatusUnauthorized = 401
)

var Need2FAErr = errors.New("need 2FA")

type Client struct {
	token []byte
}

func NewClient(token []byte) *Client {
	return &Client{token: token}
}

func Authorize(user string, password string, fingerprint string) (*Client, error) {
	return authorizeWithHeaders(user, password, fingerprint, nil)
}

func AuthorizeWith2FA(user string, password string, fingerprint string, otpCode string) (*Client, error) {
	return authorizeWithHeaders(
		user,
		password,
		fingerprint,
		map[string]string{
			"X-GitHub-OTP": otpCode,
		},
	)
}

func authorizeWithHeaders(user string, password string, fingerprint string, headers map[string]string) (*Client, error) {
	reqBody := []byte(fmt.Sprintf(`{"note":"lgtm","scopes":["repo"],"fingerprint":"%s"}`, fingerprint))

	req, err := http.NewRequest("POST", url(authURI), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	http := &http.Client{}

	resp, err := http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == httpStatusUnauthorized && resp.Header["X-Github-Otp"] != nil {
		return nil, Need2FAErr
	}

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != httpStatusCreated {
		return nil, errors.New(string(body))
	}

	var auth Authorization
	json.Unmarshal(body, &auth)

	return &Client{token: []byte(auth.Token)}, nil
}

func (c *Client) GetToken() []byte {
	if c != nil {
		return c.token
	}

	return nil
}

func (c *Client) PullRequests(repository string) ([]*PullRequest, error) {
	body, err := c.get(url(pullRequestsURI, repository))
	if err != nil {
		return nil, err
	}

	var pullRequests []*PullRequest
	json.Unmarshal([]byte(body), &pullRequests)

	return pullRequests, nil
}

func (c *Client) Reviews(repository string, pullRequestNumber int32) ([]*Review, error) {
	body, err := c.get(url(reviewsURI, repository, pullRequestNumber))
	if err != nil {
		return nil, err
	}

	var reviews []*Review
	json.Unmarshal([]byte(body), &reviews)

	return reviews, nil
}

func (c *Client) get(url string) ([]byte, error) {
	client := &http.Client{}

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

	if resp.StatusCode != httpStatusOk {
		return nil, fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}

	return body, nil
}

func url(uri string, opts ...interface{}) string {
	return fmt.Sprintf(strings.Join([]string{protocol, host, uri}, ""), opts...)
}
