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
func StartWorkingOnIssue(customIssueName string, config IssuectlConfig, issueID IssueID) error {
	profile := config.GetProfile(config.GetCurrentProfile())
	repositories := []string{}
	for _, repoName := range profile.Repositories {
		repositories = append(repositories, string(repoName))
	}

	if isIssueIdInUse(config, issueID) {
		return fmt.Errorf("issueID already found in local configuration")
	}

	Log.Infofp("🏗️", "Preparing workspace for issue %v...", issueID)

	name := string(issueID)
	if customIssueName != "" {
		name = fmt.Sprintf("%v-%v", name, customIssueName)
	}
	dirName := name
	branchName := name

	if profile.IssueBackend != "" && customIssueName == "" {
		backendConfig := config.GetBackend(profile.IssueBackend)
		issueBackend, err := getIssueBackendConfigurator(backendConfig)
		if err != nil {
			return err
		}
		generatedBranchName, err := getBranchName(config, issueBackend, profile, issueID)
		if err != nil {
			return err
		}
		branchName = generatedBranchName
	}

	issueDirPath, err := createDirectory(profile.WorkDir, dirName)
	if err != nil {
		return err
	}

	Log.Infofp("🛬", "Cloning repositories %v", repositories)

	newIssue, err := createAndAddRepositoriesToIssue(config, profile, issueID, issueDirPath, branchName, branchName, repositories)
	if err != nil {
		return err
	}

	if profile.IssueBackend != "" {
		Log.Infofp("🫡", "Marking issue as In Progress in %v", profile.IssueBackend)

		// FIXME: this is a workaround for github. we should move this to backend
		issueRepo := config.GetRepository(profile.DefaultRepository)

		backendConfig := config.GetBackend(profile.IssueBackend)
		issueBackend, err := getIssueBackendConfigurator(backendConfig)
		if err != nil {
			return err
		}
		if err := issueBackend.StartIssue(issueRepo.Owner, issueRepo.Name, issueID); err != nil {
			return err
		}
	}

	if err := config.AddIssue(newIssue); err != nil {
		return err
	}

	Log.Infofp("🚀", "Workspace for %v ready!", issueID)
	Log.Infofp("🧑‍💻", "Run `issuectl workon %v` to open it in VS Code", issueID)
	return nil
}

// isIssueIdInUse checks if issue ID is already in use
func isIssueIdInUse(config IssuectlConfig, issueID IssueID) bool {
	_, found := config.GetIssue(issueID)
	return found
}

// getIssueAndBranchName gets issue and prepares branch name
func getBranchName(config IssuectlConfig, issueBackend IssueBackend, profile *Profile, issueID IssueID) (string, error) {
	repo := config.GetRepository(profile.DefaultRepository)
	issue, err := issueBackend.GetIssue(repo.Owner, repo.Name, issueID)
	if err != nil {
		return "", fmt.Errorf(errFailedToGetIssue, err)
	}

	toReplace := []string{
		" ",
		",",
		":",
		"|",
		";",
		"(",
		")",
		"#",
		"@",
		"!",
		".",
		"$",
		"%",
		"^",
		"&",
		"*",
	}

	switch t := issue.(type) {
	default:
		return "", fmt.Errorf("Missing issue type")
	case *github.Issue:
		Log.V(5).Infof("%v", t)
		branchName := fmt.Sprintf("%v-%v", issueID, *issue.(*github.Issue).Title)
		for _, charToReplace := range toReplace {
			branchName = strings.ReplaceAll(branchName, charToReplace, "-")
		}
		return branchName, nil
	case *jira.Issue:
		Log.V(5).Infof("%v", t)
		branchName := fmt.Sprintf("%v-%v", issueID, issue.(*jira.Issue).Fields.Summary)
		for _, charToReplace := range toReplace {
			branchName = strings.ReplaceAll(branchName, charToReplace, "-")
		}
		return branchName, nil
	}
}

// createAndAddRepositoriesToIssue prepares issue and clones repositories to it
func createAndAddRepositoriesToIssue(
	config IssuectlConfig, profile *Profile, issueID IssueID, issueDirPath string, branchName, issueTitle string, repositories []string) (*IssueConfig, error) {
	newIssue := &IssueConfig{
		Name:         issueTitle,
		ID:           issueID,
		BranchName:   branchName,
		RepoBackend:  profile.RepoBackend,
		IssueBackend: profile.IssueBackend,
		Dir:          issueDirPath,
		Profile:      profile.Name,
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

	Log.V(3).Infof("Cloning repo %v", repo.Name)

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
	issue, found := config.GetIssue(issueID)
	if !found {
		return errors.New("Issue not found")
	}
	profile := config.GetProfile(issue.Profile)

	if profile.RepoBackend == "" {
		return errors.New("Repository Backend not defined")
	}

	repoBackend, err := getRepoBackendConfigurator(config.GetBackend(profile.RepoBackend))
	if err != nil {
		return err
	}

	repo := config.GetRepository(profile.DefaultRepository)

	Log.Infofp("📂", "Opening PR for issue %v in %v/%v [%v]",
		issueID,
		repo.Owner,
		repo.Name,
		profile.RepoBackend,
	)

	prId, err := repoBackend.OpenPullRequest(
		repo.Owner,
		repo.Name,
		fmt.Sprintf("%v | %v", issue.ID, issue.Name),
		fmt.Sprintf("Resolves #%v ✅", issue.ID),
		"master", // TODO: make configurable
		issue.BranchName,
	)
	if err != nil {
		return err
	}

	issueBackend, err := getIssueBackendConfigurator(config.GetBackend(profile.IssueBackend))
	if err != nil {
		return err
	}
	Log.Infofp("🔗", "Linking PR %v to issue %v in %v", *prId, issueID, profile.IssueBackend)

	return issueBackend.LinkIssueToRepo(repo.Owner, repo.Name, issueID, strconv.Itoa(*prId))
}

// FinishWorkingOnIssue finishes work on an issue
func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig().GetPersistent()
	issue, found := config.GetIssue(issueID)
	if !found {
		return errors.New("Issue not found")
	}
	profile := config.GetProfile(issue.Profile)

	repo := config.GetRepository(profile.DefaultRepository)

	Log.Infofp("🥂", "Finishing work on %v", issueID)
	if profile.IssueBackend != "" {
		issueBackend, err := getIssueBackendConfigurator(config.GetBackend(profile.IssueBackend))
		if err != nil {
			return err
		}

		Log.Infofp("🏁", "Closing issue %v in %v", issueID, profile.IssueBackend)

		err = issueBackend.CloseIssue(
			repo.Owner,
			repo.Name,
			issueID,
		)
		if err != nil {
			return fmt.Errorf(errFailedToCloseIssue, err)
		}

	}

	Log.Infofp("🧹", "Cleaning up issue workdir")

	if err := os.RemoveAll(issue.Dir); err != nil {
		return err
	}

	Log.Infofp("🫥", "Removing issue config")

	if err := config.DeleteIssue(issueID); err != nil {
		return err
	}

	Log.Infofp("👍", "All done!")

	return nil
}
