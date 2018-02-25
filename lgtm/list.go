package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Repo struct {
	PullRequests []PullRequest
	Name         string
	Error        error
}

func List(showAll bool, secrets *Secrets, config *Config) error {
	if !secrets.CheckToken() {
		fmt.Println("Please log in first (lgtm login)")
		return nil
	}

	repoCh := make(chan Repo)

	repoNames := config.Repos

	for _, repo := range repoNames {
		go listRepo(repo, config.UserName, showAll, secrets, repoCh)
	}

	reposWithPrs := make(map[string][]PullRequest)

	for range repoNames {
		repo := <-repoCh

		if repo.Error != nil {
			errWithRepo := fmt.Errorf("%s (repository: %s)", repo.Error.Error(), repo.Name)

			return errWithRepo
		}

		if len(repo.PullRequests) > 0 {
			reposWithPrs[repo.Name] = repo.PullRequests
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
		repoCh <- Repo{Error: err, Name: repo}
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
			reviews, err := getReviews(repo, pr, secrets)

			pr.Reviews = reviews

			if err != nil {
				repoCh <- Repo{Error: err, Name: repo}
				return
			}

			prs = append(prs, pr)
		}
	}

	repoCh <- Repo{Name: repo, PullRequests: prs}
}

func getReviews(repo string, pr PullRequest, secrets *Secrets) ([]Review, error) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls/%d/reviews", repo, pr.Number), secrets)

	if err != nil {
		return nil, err
	}

	var reviews []Review
	json.Unmarshal([]byte(body), &reviews)

	return reviews, nil
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

	out := [][]string{}

	for repo, prs := range repos {
		for _, pr := range prs {
			states := []string{}

			rejections := 0
			approvals := 0
			comments := 0

			for _, review := range pr.Reviews {
				states = append(states, review.State)
				switch review.State {
				case "APPROVED":
					approvals += 1
				case "CHANGES_REQUESTED":
					rejections += 1
				case "COMMENTED":
					comments += 1
				}
			}

			var state string

			if rejections > 0 {
				state = "❌"
			} else if approvals > 1 {
				state = "✅"
			} else if comments > 0 || approvals == 0 {
				state = "💬"
			}

			out = append(out, []string{repo, fmt.Sprintf("%d", pr.Number), pr.User.Login, pr.Title, state})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repository", "Number", "User", "Title", "State"})
	table.SetColWidth(100)
	table.SetBorder(false)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_CENTER,
	})
	table.AppendBulk(out)
	table.Render()
}
