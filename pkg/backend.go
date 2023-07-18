package issuectl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// IssueBackend is an interface that defines the methods that an issue backend should have.
type IssueBackend interface {
	IssueExists(owner, repo, issueID string) (bool, error)
	LinkIssueToRepo(owner, repo, issueID, pullRequestID, token string) error
	CloseIssue(owner, repo, issueID, token string) error
}

// GitHub is a struct that implements the IssueBackend interface for GitHub.
type GitHub struct{}

// IssueExists checks if an issue with the given ID exists in the specified GitHub repository.
func (g *GitHub) IssueExists(owner, repo, issueID string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", owner, repo, issueID)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// LinkIssueToRepo links a pull request to an issue in the specified GitHub repository.
func (g *GitHub) LinkIssueToRepo(owner, repo, issueID, pullRequestID, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/timeline", owner, repo, issueID)
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

// CloseIssue closes an issue in the specified GitHub repository.
func (g *GitHub) CloseIssue(owner, repo, issueID, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", owner, repo, issueID)
	body := map[string]string{
		"state": "closed",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close issue: status code %d", resp.StatusCode)
	}

	return nil
}
