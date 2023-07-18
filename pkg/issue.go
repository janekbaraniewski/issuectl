package issuectl

import (
	"fmt"
	"os"
	"path/filepath"
)

func StartWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()
	profile := config.GetProfile(config.CurrentProfile)
	// backendConfig := config.GetBackend(profile.Backend)
	if existing := config.GetIssue(issueID); existing != nil {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)
	Log.V(2).Infof("Creating issue work dir")
	issueDirPath, err := createDirectory(profile.WorkDir, string(issueID))
	if err != nil {
		return err
	}
	Log.V(2).Infof("Cloning repo")
	repoDirPath, err := cloneRepo(&config.Repositories[0], issueDirPath) // FIXME: should use repo name to get repo config instead of getting it with direct array access
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

func FinishWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()

	Log.Infof("Cleaning up after work on issue %v", issueID)
	if err := os.RemoveAll(filepath.Join(config.WorkDir, string(issueID))); err != nil {
		return err
	}

	return config.DeleteIssue(issueID)
}
