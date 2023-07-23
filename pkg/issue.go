package issuectl

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/google/go-github/github"
)

const (
	GitHubApi                     = "https://api.github.com/"
	errIssueIDNotFound            = "issueID not found"
	errIssueDoesNotExistOnBackend = "while looking for issue on %v: %v"
	errFailedToCloseIssue         = "failed to close the issue: %w"
	errFailedToGetIssue           = "failed to get the issue: %w"
)

// StartWorkingOnIssue starts work on an issue
func StartWorkingOnIssue(config IssuectlConfig, issueID IssueID) error {
	profile := config.GetProfile(config.GetCurrentProfile())
	repositories := []string{}
	for _, repoName := range profile.Repositories {
		Log.Infof("Appending repo %v", repoName)
		repositories = append(repositories, string(repoName))
	}

	if isIssueIdInUse(config, issueID) {
		return fmt.Errorf("issueID already found in local configuration")
	}

	Log.Infof("Starting work on issue %v ...", issueID)

	issueBackend, issueDirPath, err := initializeIssueBackendAndDir(config, profile, issueID)
	if err != nil {
		return err
	}

	branchName, err := getBranchName(config, issueBackend, profile, issueID)
	if err != nil {
		return err
	}

	newIssue, err := createAndAddRepositoriesToIssue(config, profile, issueID, issueDirPath, branchName, branchName, repositories)
	if err != nil {
		return err
	}

	if err := issueBackend.StartIssue("", "", issueID); err != nil {
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
	backendConfig := config.GetBackend(profile.IssueBackend)
	issueBackend, err := getIssueBackendConfigurator(backendConfig)
	if err != nil {
		return nil, "", err
	}

	repo := config.GetRepository(profile.DefaultRepository)
	issueFromBackend, err := issueBackend.GetIssue(repo.Owner, repo.Name, issueID)
	if err != nil {
		return nil, "", fmt.Errorf(errIssueDoesNotExistOnBackend, backendConfig.Name, err)
	}
	if issueFromBackend == nil {
		return nil, "", fmt.Errorf("Issue %v not found in backend %v", issueID, backendConfig.Name)
	}

	issueDirPath, err := createDirectory(profile.WorkDir, string(issueID))
	if err != nil {
		return nil, "", err
	}

	return issueBackend, issueDirPath, nil
}

// getIssueAndBranchName gets issue and prepares branch name
func getBranchName(config IssuectlConfig, issueBackend IssueBackend, profile *Profile, issueID IssueID) (string, error) {
	repo := config.GetRepository(profile.DefaultRepository)
	issue, err := issueBackend.GetIssue(repo.Owner, repo.Name, issueID)
	if err != nil {
		return "", fmt.Errorf(errFailedToGetIssue, err)
	}

	switch t := issue.(type) {
	default:
		return "", fmt.Errorf("Missing issue type")
	case *github.Issue:
		Log.Infof("%v", t)
		branchName := fmt.Sprintf("%v-%v", issueID, strings.ReplaceAll(*issue.(*github.Issue).Title, " ", "-"))
		return branchName, nil
	case *jira.Issue:
		Log.Infof("%v", t)
		branchName := fmt.Sprintf("%v-%v", issueID, strings.ReplaceAll(issue.(*jira.Issue).Fields.Summary, " ", "-"))
		return branchName, nil
	}
}

// createAndAddRepositoriesToIssue prepares issue and clones repositories to it
func createAndAddRepositoriesToIssue(
	config IssuectlConfig, profile *Profile, issueID IssueID, issueDirPath string, branchName, issueTitle string, repositories []string) (*IssueConfig, error) {
	newIssue := &IssueConfig{
		Name:        issueTitle,
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
	if repo == nil {
		return fmt.Errorf("Repo %v not defined.", repoName)
	}

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

	repoBackend, err := getRepoBackendConfigurator(config.GetBackend(profile.RepoBackend))
	if err != nil {
		return err
	}

	repo := config.GetRepository(profile.DefaultRepository)
	prId, err := repoBackend.OpenPullRequest(
		repo.Owner,
		repo.Name,
		fmt.Sprintf("%v | %v", issue.ID, issue.Name),
		fmt.Sprintf("Resolves #%v", issue.ID),
		"master", // TODO: make configurable
		issue.BranchName,
	)
	if err != nil {
		return err
	}

	issueBackend, err := getIssueBackendConfigurator(config.GetBackend(profile.RepoBackend))
	if err != nil {
		return err
	}
	return issueBackend.LinkIssueToRepo(repo.Owner, repo.Name, issueID, strconv.Itoa(*prId))
}

// FinishWorkingOnIssue finishes work on an issue
func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig().GetPersistent()
	profile := config.GetProfile(config.GetCurrentProfile())
	repo := config.GetRepository(profile.DefaultRepository)

	issueBackend, err := getIssueBackendConfigurator(config.GetBackend(profile.IssueBackend))
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
