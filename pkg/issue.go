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
	GitHubApi                    = "https://api.github.com/"
	errIssueIDNotFound           = "issueID not found"
	errIssueDoesNotExistOnGitHub = "issue does not exist on GitHub: %w"
	errFailedToCloseIssue        = "failed to close the issue: %w"
	errFailedToGetIssue          = "failed to get the issue: %w"
)

// decodeBackendToken decodes backend token from base64
func decodeBackendToken(backendConfig *BackendConfig) (string, error) {
	ghToken, err := base64.RawStdEncoding.DecodeString(backendConfig.Token)
	if err != nil {
		return "", err
	}
	return string(ghToken), nil
}

// getIssueBackendConfigurator prepares IssueBackend
func getIssueBackendConfigurator(backendConfig *BackendConfig) (IssueBackend, error) {
	ghToken, err := decodeBackendToken(backendConfig)
	if err != nil {
		return nil, err
	}
	return GetIssueBackend(&GetBackendConfig{
		Type:        backendConfig.Type,
		GitHubApi:   GitHubApi,
		GitHubToken: ghToken,
	}), nil
}

// getRepoBackendConfigurator prepares RepositoryBackend
func getRepoBackendConfigurator(backendConfig *BackendConfig) (RepositoryBackend, error) {
	ghToken, err := decodeBackendToken(backendConfig)
	if err != nil {
		return nil, err
	}
	return GetRepoBackend(&GetBackendConfig{
		Type:        backendConfig.Type,
		GitHubApi:   GitHubApi,
		GitHubToken: ghToken,
	}), nil
}

// StartWorkingOnIssue starts work on an issue
func StartWorkingOnIssue(config IssuectlConfig, overwrites *CLIOverwrites, issueID IssueID) error {
	profile := config.GetProfile(config.GetCurrentProfile())
	repositories := []string{}
	for _, repoName := range profile.Repositories {
		repositories = append(repositories, string(repoName))
	}

	if isIssueIdInUse(config, issueID) {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)

	issueBackend, issueDirPath, err := initializeIssueBackendAndDir(config, profile, issueID)
	if err != nil {
		return err
	}

	issue, branchName, err := getIssueAndBranchName(config, issueBackend, profile, issueID)
	if err != nil {
		return err
	}

	newIssue, err := createAndAddRepositoriesToIssue(config, profile, issueID, issueDirPath, branchName, issue, repositories)
	if err != nil {
		return err
	}

	if err := config.AddIssue(newIssue); err != nil {
		return err
	}

	Log.Infof("Started working on issue %v", issueID)
	return nil
}

// isIssueIdInUse checks if issue ID is already in use
func isIssueIdInUse(config IssuectlConfig, issueID IssueID) bool {
	_, found := config.GetIssue(issueID)
	return found
}

// initializeIssueBackendAndDir prepares IssueBackend and creates directory for issue
func initializeIssueBackendAndDir(config IssuectlConfig, profile *Profile, issueID IssueID) (IssueBackend, string, error) {
	backendConfig := config.GetBackend(profile.Backend)
	issueBackend, err := getIssueBackendConfigurator(backendConfig)
	if err != nil {
		return nil, "", err
	}

	repo := config.GetRepository(profile.DefaultRepository)
	exists, err := issueBackend.IssueExists(repo.Owner, repo.Name, issueID)
	if err != nil || !exists {
		return nil, "", fmt.Errorf(errIssueDoesNotExistOnGitHub, err)
	}

	issueDirPath, err := createDirectory(profile.WorkDir, string(issueID))
	if err != nil {
		return nil, "", err
	}

	return issueBackend, issueDirPath, nil
}

// getIssueAndBranchName gets issue and prepares branch name
func getIssueAndBranchName(config IssuectlConfig, issueBackend IssueBackend, profile *Profile, issueID IssueID) (interface{}, string, error) {
	repo := config.GetRepository(profile.DefaultRepository)
	issue, err := issueBackend.GetIssue(repo.Owner, repo.Name, issueID)
	if err != nil {
		return nil, "", fmt.Errorf(errFailedToGetIssue, err)
	}

	branchName := fmt.Sprintf("%v-%v", issueID, strings.ReplaceAll(*issue.(*github.Issue).Title, " ", "-"))
	return issue, branchName, nil
}

// createAndAddRepositoriesToIssue prepares issue and clones repositories to it
func createAndAddRepositoriesToIssue(config IssuectlConfig, profile *Profile, issueID IssueID, issueDirPath string, branchName string, issue interface{}, repositories []string) (*IssueConfig, error) {
	newIssue := &IssueConfig{
		Name:        *issue.(*github.Issue).Title,
		ID:          issueID,
		BranchName:  branchName,
		BackendName: "github",
		Dir:         issueDirPath,
		Profile:     profile.Name,
	}

	for _, repoName := range repositories {
		err := cloneAndAddRepositoryToIssue(config, profile, newIssue, issueDirPath, branchName, repoName)
		if err != nil {
			return nil, err
		}
	}

	return newIssue, nil
}

// cloneAndAddRepositoryToIssue clones repository and adds it to issue
func cloneAndAddRepositoryToIssue(config IssuectlConfig, profile *Profile, issue *IssueConfig, issueDirPath string, branchName string, repoName string) error {
	gitUser, _ := config.GetGitUser(profile.GitUserName)
	repo := config.GetRepository(RepoConfigName(repoName))

	Log.Infof("Cloning repo %v", repo.Name)

	repoDirPath, err := cloneRepo(repo, issueDirPath, gitUser)
	if err != nil {
		return err
	}

	Log.V(2).Infof("Creating branch")
	if err := createBranch(repoDirPath, branchName, gitUser); err != nil {
		return err
	}

	issue.Repositories = append(issue.Repositories, repo.Name)
	return nil
}

// OpenPullRequest opens pull request
func OpenPullRequest(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.GetCurrentProfile())

	issue, found := config.GetIssue(issueID)
	if !found {
		return fmt.Errorf(errIssueIDNotFound)
	}

	repoBackend, err := getRepoBackendConfigurator(config.GetBackend(profile.Backend))
	if err != nil {
		return err
	}

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

// FinishWorkingOnIssue finishes work on an issue
func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.GetCurrentProfile())
	repo := config.GetRepository(profile.DefaultRepository)

	issueBackend, err := getIssueBackendConfigurator(config.GetBackend(profile.Backend))
	if err != nil {
		return err
	}

	err = issueBackend.CloseIssue(
		repo.Owner,
		repo.Name,
		issueID,
	)
	if err != nil {
		return fmt.Errorf(errFailedToCloseIssue, err)
	}

	Log.Infof("Cleaning up after work on issue %v", issueID)

	issue, found := config.GetIssue(issueID)
	if !found {
		return errors.New(errIssueIDNotFound)
	}

	if err := os.RemoveAll(issue.Dir); err != nil {
		return err
	}

	return config.DeleteIssue(issueID)
}
