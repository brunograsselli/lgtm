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
	prsCh := make(chan []PullRequest)

	for _, repo := range repos {
		go listRepo(repo, user, token, prsCh)
	}

	reposWithPrs := make(map[string][]PullRequest)

	for _, repo := range repos {
		repoPrs := <-prsCh

		if len(repoPrs) > 0 {
			reposWithPrs[repo] = repoPrs
		}
	}

	if len(reposWithPrs) == 0 {
		fmt.Println("You are up to date!")
		return
	}

	for repo, prs := range reposWithPrs {
		fmt.Printf("%s:\n", repo)
		for _, pr := range prs {
			fmt.Printf("  %s\n", pr.Title)
		}
		fmt.Println("")
	}
}

func listRepo(repo string, user string, token string, prsCh chan []PullRequest) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo), nil)

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

	prs := []PullRequest{}

	var includePR bool

	for _, pr := range pullRequests {
		includePR = false

		for _, reviewer := range pr.RequestedReviewers {
			if reviewer.Login == user {
				includePR = true
			}
		}

		if includePR {
			prs = append(prs, pr)
		}
	}

	prsCh <- prs
}
