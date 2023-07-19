package cli

import (
	"fmt"
	"os"
	"os/exec"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"

	"github.com/spf13/cobra"
)

func initOpenIssueCommand(rootCmd *cobra.Command) {
	var openIssueCmd = &cobra.Command{
		Use:   "open [issueID]",
		Short: "Open specified issue in the preferred code editor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig() // Load your configuration here
			issueID := issuectl.IssueID(args[0])

			issue, found := config.GetIssue(issueID)
			if !found {
				return fmt.Errorf("Issue %s not found", issueID)
			}

			// Open the preferred editor with the directory IssueConfig.Dir
			openCmd := exec.Command("code", issue.Dir) // Change 'vim' to your preferred editor
			openCmd.Stdin = os.Stdin
			openCmd.Stdout = os.Stdout
			err := openCmd.Run()
			if err != nil {
				return fmt.Errorf("Failed to open editor: %v", err)
			}

			return nil
		},
	}

	rootCmd.AddCommand(openIssueCmd)
}
