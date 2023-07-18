package cli

import (
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initOpenPullRequestCommand(rootCmd *cobra.Command) {
	openPRCmd := &cobra.Command{
		Use:   "openpr",
		Short: "Opens a pull request for the specified issue",
		Long:  "This command opens a pull request for the specified issue on GitHub",
		Args:  cobra.ExactArgs(1), // it requires exactly one argument
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0] // the first argument is the issue ID
			err := issuectl.OpenPullRequest(issuectl.IssueID(issueID))
			if err != nil {
				issuectl.Log.Infof("Error!! -> %v", err)
				return err
			}
			return nil
		},
	}

	rootCmd.AddCommand(openPRCmd)
}
