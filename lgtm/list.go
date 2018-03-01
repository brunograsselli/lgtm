package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type Repo struct {
	PullRequests []PullRequest
	Name         string
	Error        error
}

type ReviewsWithError struct {
	Reviews           []Review
	Error             error
	PullRequestNumber int32
}

func List(showAll bool, secrets *Secrets, config *Config) error {
	if !secrets.CheckToken() {
		err := Login(secrets, config)

		if err != nil {
			return err
		}
	}

	repoNames := config.Repos

	repoCh := make(chan Repo, len(repoNames))
	var wg sync.WaitGroup

	for _, name := range repoNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			prs, err := fetchRepo(name, config.UserName, showAll, secrets)
			repoCh <- Repo{Error: err, Name: name, PullRequests: prs}
		}(name)
	}

	go func() {
		wg.Wait()
		close(repoCh)
	}()

	reposWithPrs := make(map[string][]PullRequest)

	for repo := range repoCh {
		if repo.Error != nil {
			return fmt.Errorf("%s (repository: %s)", repo.Error.Error(), repo.Name)
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

func fetchRepo(repo string, user string, showAll bool, secrets *Secrets) ([]PullRequest, error) {
	openPrs, err := GitHubPullRequests(repo, secrets)

	if err != nil {
		return nil, err
	}

	filteredPrs := []PullRequest{}

	var includePR bool
	var wg sync.WaitGroup
	reviewsCh := make(chan ReviewsWithError, 5)

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
			wg.Add(1)

			go func(number int32) {
				defer wg.Done()
				reviews, err := GitHubReviews(repo, number, secrets)

				reviewsCh <- ReviewsWithError{Error: err, Reviews: reviews, PullRequestNumber: number}
			}(pr.Number)

			filteredPrs = append(filteredPrs, pr)
		}
	}

	go func() {
		wg.Wait()
		close(reviewsCh)
	}()

	reviews := make(map[int32][]Review)

	for r := range reviewsCh {
		if r.Error != nil {
			return nil, r.Error
		}

		reviews[r.PullRequestNumber] = r.Reviews
	}

	prsWithReviews := []PullRequest{}

	for _, pr := range filteredPrs {
		pr.Reviews = reviews[pr.Number]
		prsWithReviews = append(prsWithReviews, pr)
	}

	return prsWithReviews, nil
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
