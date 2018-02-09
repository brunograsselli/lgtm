package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func List(showAll bool) error {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		fmt.Println("Please log in first (lgtm login)")
		return nil
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

	err := write(reposWithPrs)
	if err != nil {
		return err
	}
	print(reposWithPrs)
	return nil
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

func write(repos map[string][]PullRequest) error {
	c, err := json.Marshal(repos)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/tmp/lgtm.json", c, 0644)

	if err != nil {
		return err
	}

	return nil
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
