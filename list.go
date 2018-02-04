package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

type User struct {
	Login string `json:"login"`
}

type PullRequest struct {
	URL                string `json:"url"`
	Number             int32  `json:"number"`
	Title              string `json:"title"`
	User               User   `json:"user"`
	RequestedReviewers []User `json:"requested_reviewers"`
}

func List() {
	repos := viper.GetStringSlice("repos")
	user := viper.GetString("user")
	token := viper.GetString("token")

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls", repos[0]), nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}

	var pullRequests []PullRequest
	json.Unmarshal([]byte(body), &pullRequests)

	pendingPRs := FilterPullRequest(pullRequests, func(p PullRequest) bool {
		return IncludeString(MapUserString(p.RequestedReviewers, func(user User) string {
			return user.Login
		}), user)
	})

	fmt.Println(pendingPRs)
}
