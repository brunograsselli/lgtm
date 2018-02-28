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
		err := Login(secrets, config)

		if err != nil {
			return err
		}
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

	err := writeTempFile(reposWithPrs)

	if err == nil {
		printList(reposWithPrs)
	}

	return err
}

func listRepo(repo string, user string, showAll bool, secrets *Secrets, repoCh chan Repo) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls?sort=created&direction=asc", repo), secrets)

	if err != nil {
		repoCh <- Repo{Error: err, Name: repo}
		return
	}

	var openPrs []PullRequest
	json.Unmarshal([]byte(body), &openPrs)

	filteredPrs := []PullRequest{}

	var includePR bool

	for _, pr := range openPrs {
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

			filteredPrs = append(filteredPrs, pr)
		}
	}

	repoCh <- Repo{Name: repo, PullRequests: filteredPrs}
}

func getReviews(repo string, pr PullRequest, secrets *Secrets) ([]Review, error) {
	body, err := GitHubGet(fmt.Sprintf("/repos/%s/pulls/%d/reviews?sort=created&direction=asc", repo, pr.Number), secrets)

	if err != nil {
		return nil, err
	}

	var reviews []Review
	json.Unmarshal([]byte(body), &reviews)

	return reviews, nil
}

func writeTempFile(repos map[string][]PullRequest) error {
	c, err := json.Marshal(repos)

	if err != nil {
		return err
	}

	return ioutil.WriteFile("/tmp/lgtm.json", c, 0644)
}

func printList(repos map[string][]PullRequest) {
	if len(repos) == 0 {
		fmt.Println("You are up to date!")
		return
	}

	out := [][]string{}

	for repo, prs := range repos {
		for _, pr := range prs {
			state := computeState(pr)

			out = append(out, []string{repo, fmt.Sprint(pr.Number), pr.User.Login, pr.Title, state})
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

func computeState(pr PullRequest) string {
	states := make(map[string]string)

	// Consider only the last review of each user
	for _, review := range pr.Reviews {
		states[review.User.Login] = review.State
	}

	var rejections, approvals, comments int

	for _, state := range states {
		switch state {
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
	} else if comments > 0 || approvals == 1 {
		state = "💬"
	}

	return state
}
