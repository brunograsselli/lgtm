package lgtm

import (
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
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}

	return body, nil
}