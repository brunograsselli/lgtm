package github

type PullRequest struct {
	URL                string  `json:"url"`
	HTMLURL            string  `json:"html_url"`
	Number             int32   `json:"number"`
	Title              string  `json:"title"`
	User               User    `json:"user"`
	RequestedReviewers []*User `json:"requested_reviewers"`
	Reviews            []*Review
}

type Review struct {
	User  User   `json:"user"`
	State string `json:"state"`
}

type User struct {
	Login string `json:"login"`
}

type Authorization struct {
	Token string `json:"token"`
}
