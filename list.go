package lgtm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

func List(showAll bool) {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		fmt.Println("Please log in first (lgtm login)")
		return
	}

	repos := viper.GetStringSlice("repos")
	user := viper.GetString("user")
	prsCh := make(chan []PullRequest)

	for _, repo := range repos {
		go listRepo(repo, user, showAll, prsCh)
	}

	reposWithPrs := make(map[string][]PullRequest)

	for _, repo := range repos {
		repoPrs := <-prsCh

		if len(repoPrs) > 0 {
			reposWithPrs[repo] = repoPrs
		}
	}

	print(reposWithPrs)
}

func listRepo(repo string, user string, showAll bool, prsCh chan []PullRequest) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls", repo))

	if err != nil {
		panic(err)
	}

	var pullRequests []PullRequest
	json.Unmarshal([]byte(body), &pullRequests)

	prs := []PullRequest{}

	var includePR bool

	for _, pr := range pullRequests {
		includePR = false

		if showAll {
			includePR = true
		} else {
			for _, reviewer := range pr.RequestedReviewers {
				if reviewer.Login == user {
					includePR = true
				}
			}
		}

		if includePR {
			prs = append(prs, pr)
		}
	}

	prsCh <- prs
}

func print(repos map[string][]PullRequest) {
	if len(repos) == 0 {
		fmt.Println("You are up to date!")
		return
	}

	out := []string{}

	for repo, prs := range repos {
		out = append(out, fmt.Sprintf("%s:", repo))

		for _, pr := range prs {
			out = append(out, fmt.Sprintf("  %d\t%s", pr.Number, pr.Title))
		}

		out = append(out, "")
	}

	fmt.Println(strings.Join(out[0:len(out)-1], "\n"))
}
