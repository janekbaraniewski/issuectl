package issuectl

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

const (
	GitHubApi = "https://api.github.com/"
)

func StartWorkingOnIssue(config *IssuectlConfig, repositories []string, issueID IssueID) error {
	profile := config.GetProfile(config.GetCurrentProfile())
	for _, repoName := range repositories {
		// TODO: make sure that repos are not duplicated OR overwrite instead of adding repositories
		if err := profile.AddRepository((*RepoConfigName)(&repoName)); err != nil {
			return err
		}
	}
	if _, found := config.GetIssue(issueID); found {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)

	repo := config.GetRepository(profile.DefaultRepository)
	backendConfig := config.GetBackend(profile.Backend)
	gitUser, _ := config.GetGitUser(profile.GitUserName)
	ghToken, err := base64.RawStdEncoding.DecodeString(backendConfig.Token)
	if err != nil {
		return err
	}
	issueBackend := GetIssueBackend(&GetBackendConfig{
		Type:        backendConfig.Type,
		GitHubApi:   GitHubApi,
		GitHubToken: string(ghToken),
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

	Log.Infof("Cloning multiple repositories: %v", profile.Repositories)
	newIssue := &IssueConfig{
		Name:        *issue.(*github.Issue).Title,
		ID:          issueID,
		BranchName:  branchName,
		BackendName: "github",
		Dir:         issueDirPath,
		Profile:     profile.Name,
	}
	for _, repoName := range profile.Repositories {
		repo := config.GetRepository(*repoName)
		Log.Infof("Cloning repo %v", repo.Name)
		repoDirPath, err := cloneRepo(&repo, issueDirPath, &gitUser)
		if err != nil {
			return err
		}
		Log.V(2).Infof("Creating branch")
		if err := createBranch(repoDirPath, branchName, &gitUser); err != nil {
			return err
		}
		newIssue.Repositories = append(newIssue.Repositories, repo.Name)
	}
	if err := config.AddIssue(newIssue); err != nil {
		return err
	}

	Log.Infof("Started working on issue %v", issueID)
	return nil
}

func OpenPullRequest(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.GetCurrentProfile())

	issue, found := config.GetIssue(issueID)
	if !found {
		return fmt.Errorf("issueID not found")
	}

	backendConfig := config.GetBackend(profile.Backend)
	ghToken, err := base64.RawStdEncoding.DecodeString(backendConfig.Token)
	if err != nil {
		return err
	}
	repoBackend := GetRepoBackend(&GetBackendConfig{
		Type:        backendConfig.Type,
		GitHubApi:   GitHubApi,
		GitHubToken: string(ghToken),
	})
	repo := config.GetRepository(profile.DefaultRepository)
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
	repo := config.GetRepository(profile.DefaultRepository)
	backendConfig := config.GetBackend(profile.Backend)
	ghToken, err := base64.RawStdEncoding.DecodeString(backendConfig.Token)
	if err != nil {
		return err
	}
	issueBackend := GetIssueBackend(&GetBackendConfig{
		Type:        backendConfig.Type,
		GitHubToken: string(ghToken),
		GitHubApi:   GitHubApi,
	})
	err = issueBackend.CloseIssue(
		repo.Owner,
		repo.Name,
		issueID,
	)
	if err != nil {
		return fmt.Errorf("failed to close the issue: %v", err)
	}

	Log.Infof("Cleaning up after work on issue %v", issueID)

	issue, found := config.GetIssue(issueID)
	if !found {
		return errors.New("Issue not found")
	}

	if err := os.RemoveAll(issue.Dir); err != nil {
		return err
	}

	return config.DeleteIssue(issueID)
}
