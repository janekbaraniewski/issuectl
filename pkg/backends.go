package issuectl

import (
	"fmt"
	"strconv"
)

// getIssueNumberFromString converts IssueID to int
func getIssueNumberFromString(issueID IssueID) (int, error) {
	issueNumber, err := strconv.Atoi(string(issueID))
	if err != nil {
		return 0, fmt.Errorf("invalid issue ID: %v", err)
	}

	return issueNumber, nil
}

type IssueBackend interface {
	IssueExists(owner string, repo RepoConfigName, issueID IssueID) (bool, error)
	LinkIssueToRepo(owner string, repo RepoConfigName, issueID IssueID, pullRequestID string) error
	CloseIssue(owner string, repo RepoConfigName, issueID IssueID) error
	GetIssue(owner string, repo RepoConfigName, issueID IssueID) (interface{}, error)
}

type RepositoryBackend interface {
	OpenPullRequest(owner string, repo RepoConfigName, title, body, baseBranch, headBranch string) error
}

type GetBackendConfig struct {
	Type        BackendType
	GitHubToken string
	GitHubApi   string
	GitLabToken string
	GitLabApi   string
}

func GetIssueBackend(conf *GetBackendConfig) IssueBackend {
	switch conf.Type {
	case BackendGithub:
		return NewGitHubClient(
			conf.GitHubToken,
			conf.GitHubApi,
		)
	case BackendGitLab:
		return NewGitLabClient(
			conf.GitLabToken,
			conf.GitLabApi,
		)
	}
	return nil
}

func GetRepoBackend(conf *GetBackendConfig) RepositoryBackend {
	switch conf.Type {
	case BackendGithub:
		return NewGitHubClient(
			conf.GitHubToken,
			conf.GitHubApi,
		)
	case BackendGitLab:
		return NewGitLabClient(
			conf.GitLabToken,
			conf.GitLabApi,
		)
	}
	return nil
}
