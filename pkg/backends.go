package issuectl

import (
	"fmt"
	"strconv"
)

// getIssueNumberFromString converts IssueID to int
func getIssueNumberFromString(issueID IssueID) (int, error) {
	issueNumber, err := strconv.Atoi(string(issueID))
	if err != nil {
		return 0, fmt.Errorf("issueID has to be of type int")
	}

	return issueNumber, nil
}

type IssueBackend interface {
	LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error
	CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error
	GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error)

	// Deprecated
	IssueExists(owner string, repo RepoConfigName, issueID IssueID) (bool, error)
}

type RepositoryBackend interface {
	OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) error
}
