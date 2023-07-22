package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func isSubdir(baseDir, checkDir string) (string, error) {
	// Get the relative path from baseDir to checkDir
	relativePath, err := filepath.Rel(baseDir, checkDir)
	if err != nil {
		return "", err
	}

	// Check if the relative path starts with ".."
	isSubdir := !strings.HasPrefix(relativePath, "..")
	if !isSubdir {
		return "", nil
	}
	return relativePath, nil
}

func getIssueIDFromParentDirectory(config issuectl.IssuectlConfig) string {
	profile := config.GetProfile(config.GetCurrentProfile())
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	isSub, err := isSubdir(profile.WorkDir, dir)
	if err != nil || isSub == "" {
		return ""
	}
	issueID := strings.Split(isSub, "/")[0]
	return issueID
}

func initOpenPullRequestCommand(rootCmd *cobra.Command) {
	openPRCmd := &cobra.Command{
		Use:   "openpr [issue number]",
		Short: "Opens a pull request for the specified issue",
		Long:  `This command opens a pull request for the specified issue on GitHub`, // FIXME: Github??                                                 // it requires exactly one argument
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			var issueID string
			if len(args) == 1 {
				issueID = args[0]
			} else if len(args) == 0 {
				// check if can detect issue from pwd
				issueID = getIssueIDFromParentDirectory(config)
			}
			if issueID == "" {
				return errors.New("Missing issueID")
			}
			err := issuectl.OpenPullRequest(issuectl.IssueID(issueID))
			if err != nil {
				return err
			}
			return nil
		},
	}

	rootCmd.AddCommand(openPRCmd)
}
