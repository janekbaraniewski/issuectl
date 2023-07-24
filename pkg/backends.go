package issuectl

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

const (
	DefaultStartMessage  = "On it 👀"
	DefaultCloseMessage  = "✅"
	DefaultOpenPRMessage = "Working on changes here: %s"
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
	StartIssue(owner string, repo RepoConfigName, issueID IssueID) error
	GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error)
	GetIssueURL(owner string, repo RepoConfigName, issueID IssueID) (string, error)
}

type RepositoryBackend interface {
	OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) (*int, error)
}

// getIssueBackendConfigurator prepares IssueBackend
func getIssueBackendConfigurator(backendConfig *BackendConfig) (IssueBackend, error) {
	switch backendConfig.Type {

	case BackendGithub:
		token, err := base64.RawStdEncoding.DecodeString(backendConfig.GitHub.Token)
		if err != nil {
			return nil, err
		}
		return NewGitHubClient(
			string(token),
			backendConfig.GitHub.Host,
			backendConfig.GitHub.Username,
		), nil

	case BackendGitLab:
		token, err := base64.RawStdEncoding.DecodeString(backendConfig.GitLab.Token)
		if err != nil {
			return nil, err
		}
		return NewGitLabClient(
			string(token),
			backendConfig.GitLab.Host,
			backendConfig.GitLab.UserID,
		), nil

	case BackendJira:
		token, err := base64.RawStdEncoding.DecodeString(backendConfig.Jira.Token)
		if err != nil {
			return nil, err
		}
		return NewJiraClient(
			backendConfig.Jira.Username,
			string(token),
			backendConfig.Jira.Host,
		), nil
	default:
		return nil, fmt.Errorf("Backend %v not supported", backendConfig.Type)
	}
}

// getRepoBackendConfigurator prepares RepositoryBackend
func getRepoBackendConfigurator(backendConfig *BackendConfig) (RepositoryBackend, error) {
	switch backendConfig.Type {
	case BackendGithub:
		token, err := base64.RawStdEncoding.DecodeString(backendConfig.GitHub.Token)
		if err != nil {
			return nil, err
		}
		return NewGitHubClient(
			string(token),
			backendConfig.GitHub.Host,
			backendConfig.GitHub.Username,
		), nil

	case BackendGitLab:
		token, err := base64.RawStdEncoding.DecodeString(backendConfig.GitLab.Token)
		if err != nil {
			return nil, err
		}
		return NewGitLabClient(
			string(token),
			backendConfig.GitLab.Host,
			backendConfig.GitLab.UserID,
		), nil
	default:
		return nil, fmt.Errorf("Backend %v not supported", backendConfig.Type)
	}
}
