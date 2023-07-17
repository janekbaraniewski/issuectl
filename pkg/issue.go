package issuectl

import (
	"fmt"
	"os"
	"path/filepath"
)

func AddIssue(issueConfig *IssueConfig) error {
	config := LoadConfig()

	config.Issues = append(config.Issues, *issueConfig)

	return config.Save()
}

func DeleteIssue(issueID IssueID) error {
	config := LoadConfig()

	for i, ic := range config.Issues {
		if ic.ID == issueID {
			if i < len(config.Issues) {
				config.Issues = append(config.Issues[:i], config.Issues[i+1:]...)
			} else {
				config.Issues = config.Issues[:i]
			}
			config.Save()
			return nil
		}
	}

	return nil
}

func GetIssue(issueID IssueID) *IssueConfig {
	config := LoadConfig()

	for _, ic := range config.Issues {
		if ic.ID == issueID {
			return &ic
		}
	}

	return nil
}

func StartWorkingOnIssue(issueID IssueID) error {
	config := LoadConfig()

	if existing := GetIssue(issueID); existing != nil {
		return fmt.Errorf("issueID already in use")
	}

	Log.Infof("Starting work on issue %v ...", issueID)
	Log.V(2).Infof("Creating issue work dir")
	issueDirPath, err := createDirectory(config.WorkDir, string(issueID))
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

	if err := AddIssue(&IssueConfig{
		Name:        string(issueID),
		ID:          issueID,
		RepoName:    config.DefaultRepository,
		BackendName: "github",
		Dir:         issueDirPath,
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

	return DeleteIssue(issueID)
}
