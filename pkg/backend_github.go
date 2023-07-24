package issuectl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHub struct {
	baseURL string
	token   string
	user    string
	client  *github.Client
}

func NewGitHubClient(token, baseURL, user string) *GitHub {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GitHub{client: client, baseURL: baseURL, token: token}
}

func (g *GitHub) GetIssueURL(owner string, repo RepoConfigName, issueID IssueID) (string, error) {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, issueNumber), nil
}

func (g *GitHub) GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error) {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return nil, err
	}

	issue, _, err := g.client.Issues.Get(context.Background(), owner, string(repo), issueNumber)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (g *GitHub) CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return err
	}

	issueRequest := &github.IssueRequest{State: github.String("closed")}
	_, _, err = g.client.Issues.Edit(context.Background(), owner, string(repo), issueNumber, issueRequest)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHub) OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) (*int, error) {
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(headBranch),
		Base:                github.String(baseBranch),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := g.client.PullRequests.Create(context.Background(), owner, string(repo), newPR)
	if err != nil {
		return nil, err
	}

	return pr.Number, nil
}

func (g *GitHub) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error {
	// Convert the pull request ID from string to int
	pullRequestNumber, err := strconv.Atoi(pullRequestID)
	if err != nil {
		return err
	}

	// Create a comment on the pull request that references the issue
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("Resolves #%s", issueID)),
	}

	// Post the comment to the pull request
	_, _, err = g.client.Issues.CreateComment(context.Background(), owner, string(repo), pullRequestNumber, comment)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHub) StartIssue(owner string, repo RepoConfigName, issueID IssueID) error {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return err
	}

	issue, _, err := g.client.Issues.Get(context.Background(), owner, string(repo), issueNumber)
	if err != nil {
		return err
	}

	// Check if the issue is already labeled "In Progress"
	for _, label := range issue.Labels {
		if *label.Name == "In Progress" {
			return nil
		}
	}

	// If not, add the "In Progress" label
	labels := []string{"In Progress"}
	_, _, err = g.client.Issues.AddLabelsToIssue(context.Background(), owner, string(repo), issueNumber, labels)
	if err != nil {
		return err
	}

	// Check if the issue is already assigned to the specified user
	for _, user := range issue.Assignees {
		if *user.Login == g.user {
			return nil
		}
	}

	// If not, assign the issue to the specified user
	assignees := []string{g.user}
	_, _, err = g.client.Issues.AddAssignees(context.Background(), owner, string(repo), issueNumber, assignees)
	if err != nil {
		return err
	}

	return nil
}
