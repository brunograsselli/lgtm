package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

type Repo struct {
	PullRequests []PullRequest
	Error        error
}

func List(showAll bool, secrets *Secrets, repoNames []string) error {
	if !secrets.CheckToken() {
		fmt.Println("Please log in first (lgtm login)")
		return nil
	}

	user := viper.GetString("user")
	repoCh := make(chan Repo)

	for _, repo := range repoNames {
		go listRepo(repo, user, showAll, secrets, repoCh)
	}

	reposWithPrs := make(map[string][]PullRequest)

	for _, name := range repoNames {
		repo := <-repoCh

		if repo.Error != nil {
			return repo.Error
		}

		if len(repo.PullRequests) > 0 {
			reposWithPrs[name] = repo.PullRequests
		}
	}

	err := write(reposWithPrs)
	if err != nil {
		return err
	}
	print(reposWithPrs)
	return nil
}

func listRepo(repo string, user string, showAll bool, secrets *Secrets, repoCh chan Repo) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls", repo), secrets)

	if err != nil {
		repoCh <- Repo{Error: err}
		return
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

	repoCh <- Repo{PullRequests: prs}
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
