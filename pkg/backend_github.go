package issuectl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

func (g *GitHub) OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) error {
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(headBranch),
		Base:                github.String(baseBranch),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err := g.client.PullRequests.Create(context.Background(), owner, string(repo), newPR)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHub) LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%s/timeline", g.baseURL, owner, repo, issueID)
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
	req.Header.Set("Authorization", "Bearer "+g.token)
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
