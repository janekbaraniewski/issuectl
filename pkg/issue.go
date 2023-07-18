package issuectl

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	GitHubApi = "https://api.github.com/"
)

func loadGithubToken() string {
	content, err := os.ReadFile("gh-access-token")
	if err != nil {
		Log.Infof("FATAL - no gh access token found")
		return ""
	}
	return string(content)
}

var GitHubToken = loadGithubToken()

func StartWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.CurrentProfile)

	if existing := config.GetIssue(issueID); existing != nil {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)

	// Check if the issue exists on GitHub.
	ghClient := NewGitHubClient(GitHubToken)
	repo := config.GetRepository(profile.Repository)
	exists, err := ghClient.IssueExists(repo.Owner, string(repo.Name), string(issueID))
	if err != nil || !exists {
		return fmt.Errorf("issue does not exist on GitHub: %v", err)
	}
	Log.V(2).Infof("Creating issue work dir")
	issueDirPath, err := createDirectory(profile.WorkDir, string(issueID))
	if err != nil {
		return err
	}
	Log.V(2).Infof("Cloning repo")
	repoDirPath, err := cloneRepo(repo, issueDirPath)
	if err != nil {
		return err
	}

	Log.V(2).Infof("Creating branch")
	if err := createBranch(repoDirPath, string(issueID)); err != nil {
		return err
	}

	if err := config.AddIssue(&IssueConfig{
		Name:        string(issueID),
		ID:          issueID,
		RepoName:    profile.Repository,
		BackendName: "github",
		Dir:         issueDirPath,
		Profile:     profile.Name,
	}); err != nil {
		return err
	}

	Log.Infof("Started working on issue %v", issueID)
	return nil
}

func OpenPullRequest(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.CurrentProfile)

	if existing := config.GetIssue(issueID); existing != nil {
		return fmt.Errorf("issueID already in use")
	}
	ghClient := NewGitHubClient(GitHubToken)
	repo := config.GetRepository(profile.Repository)

	if err := ghClient.OpenPullRequest(
		repo.Owner,
		string(repo.Name),
		"PR title", // Replace with your PR title.
		"PR body",  // Replace with your PR body.
		"master",   // Replace with the base branch name.
		string(issueID),
	); err != nil {
		return fmt.Errorf("failed to open a pull request: %v", err)
	}

	if err := ghClient.LinkIssueToRepo(
		GitHubApi,
		repo.Owner,
		string(repo.Name),
		string(issueID),
		"PR number", // Replace with your PR number.
		GitHubToken, // Replace with your real GitHub token.
	); err != nil {
		return fmt.Errorf("failed to link the pull request to the issue: %v", err)
	}

	return nil
}

func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.CurrentProfile)
	repo := config.GetRepository(profile.Repository)
	// Close the issue on GitHub.
	ghClient := NewGitHubClient(GitHubToken)
	err := ghClient.CloseIssue(
		repo.Owner,
		string(repo.Name),
		string(issueID),
	)
	if err != nil {
		return fmt.Errorf("failed to close the issue: %v", err)
	}

	Log.Infof("Cleaning up after work on issue %v", issueID)

	if err := os.RemoveAll(filepath.Join(config.WorkDir, string(issueID))); err != nil {
		return err
	}

	return config.DeleteIssue(issueID)
}
