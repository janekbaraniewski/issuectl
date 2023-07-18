package issuectl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func getIssueNumberFromString(issueID string) (int, error) {
	issueNumber, err := strconv.Atoi(issueID)
	if err != nil {
		return 0, fmt.Errorf("invalid issue ID: %v", err)
	}

	return issueNumber, nil
}

type IssueBackend interface {
	IssueExists(owner, repo, issueID string) (bool, error)
	LinkIssueToRepo(owner, repo, issueID, pullRequestID string) error
	CloseIssue(owner, repo, issueID string) error
	OpenPullRequest(owner, repo, title, body, baseBranch, headBranch string) error
}

type GitHub struct {
	client *github.Client
}

func NewGitHubClient(token string) *GitHub {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GitHub{client: client}
}

func (g *GitHub) IssueExists(owner, repo, issueID string) (bool, error) {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return false, err
	}

	_, _, err = g.client.Issues.Get(context.Background(), owner, repo, issueNumber)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (g *GitHub) CloseIssue(owner, repo, issueID, token string) error {
	issueNumber, err := getIssueNumberFromString(issueID)
	if err != nil {
		return err
	}

	issueRequest := &github.IssueRequest{State: github.String("closed")}
	_, _, err = g.client.Issues.Edit(context.Background(), owner, repo, issueNumber, issueRequest)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHub) OpenPullRequest(owner, repo, title, body, baseBranch, headBranch, token string) error {
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(headBranch),
		Base:                github.String(baseBranch),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err := g.client.PullRequests.Create(context.Background(), owner, repo, newPR)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHub) LinkIssueToRepo(baseURL, owner, repo, issueID, pullRequestID, token string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%s/timeline", baseURL, owner, repo, issueID)
	body := map[string]string{
		"issue_number": pullRequestID,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.starfire-preview+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to link issue to pull request: status code %d", resp.StatusCode)
	}

	return nil
}
