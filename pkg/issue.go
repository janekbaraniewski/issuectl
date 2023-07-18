package issuectl

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
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

func StartWorkingOnIssue(repositories []string, issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.GetCurrentProfile())
	for _, repoName := range repositories {
		rc := config.GetRepository(RepoConfigName(repoName))
		if err := profile.AddRepository(&rc); err != nil {
			return err
		}
	}
	emptyIssue := &IssueConfig{}
	if existing := config.GetIssue(issueID); existing != *emptyIssue {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)

	repo := config.GetRepository(profile.Repository)
	issueBackend := GetIssueBackend(&GetBackendConfig{
		Type:        BackendGithub,
		GitHubApi:   GitHubApi,
		GitHubToken: GitHubToken,
	})
	exists, err := issueBackend.IssueExists(repo.Owner, repo.Name, issueID)
	if err != nil || !exists {
		return fmt.Errorf("issue does not exist on GitHub: %v", err)
	}
	Log.V(2).Infof("Creating issue work dir")
	issueDirPath, err := createDirectory(profile.WorkDir, string(issueID))
	if err != nil {
		return err
	}

	issue, err := issueBackend.GetIssue(repo.Owner, repo.Name, issueID)
	if err != nil {
		return fmt.Errorf("failed to get the issue: %v", err)
	}
	branchName := fmt.Sprintf("%v-%v", issueID, strings.ReplaceAll(*issue.(*github.Issue).Title, " ", "-"))

	if profile.Repositories != nil {
		Log.Infof("Cloning multiple repositories: %v", profile.Repositories)
		for _, repo := range profile.Repositories {
			Log.Infof("Cloning repo %v", repo.Name)
			repoDirPath, err := cloneRepo(repo, issueDirPath)
			if err != nil {
				return err
			}
			Log.V(2).Infof("Creating branch")
			if err := createBranch(repoDirPath, branchName); err != nil {
				return err
			}
		}
	} else {
		Log.V(2).Infof("Cloning repo")
		repoDirPath, err := cloneRepo(&repo, issueDirPath)
		if err != nil {
			return err
		}
		Log.V(2).Infof("Creating branch")
		if err := createBranch(repoDirPath, branchName); err != nil {
			return err
		}
	}

	if err := config.AddIssue(&IssueConfig{
		Name:        *issue.(*github.Issue).Title,
		ID:          issueID,
		RepoName:    profile.Repository,
		BranchName:  branchName,
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
	profile := config.GetProfile(config.GetCurrentProfile())

	issue := config.GetIssue(issueID)
	emptyIssue := &IssueConfig{}
	if issue == *emptyIssue {
		return fmt.Errorf("issueID not found")
	}

	repoBackend := GetRepoBackend(&GetBackendConfig{
		Type:        BackendGithub,
		GitHubApi:   GitHubApi,
		GitHubToken: GitHubToken,
	})
	repo := config.GetRepository(profile.Repository)
	return repoBackend.OpenPullRequest(
		repo.Owner,
		repo.Name,
		fmt.Sprintf("%v | %v", issue.ID, issue.Name),
		fmt.Sprintf("Resolves #%v", issue.ID),
		"master", // TODO: make configurable
		issue.BranchName,
	)
}

func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.GetCurrentProfile())
	repo := config.GetRepository(profile.Repository)

	issueBackend := GetIssueBackend(&GetBackendConfig{
		Type:        BackendGithub,
		GitHubToken: GitHubToken,
		GitHubApi:   GitHubApi,
	})
	err := issueBackend.CloseIssue(
		repo.Owner,
		repo.Name,
		issueID,
	)
	if err != nil {
		return fmt.Errorf("failed to close the issue: %v", err)
	}

	Log.Infof("Cleaning up after work on issue %v", issueID)

	issue := config.GetIssue(issueID)

	if err := os.RemoveAll(issue.Dir); err != nil {
		return err
	}

	return config.DeleteIssue(issueID)
}
