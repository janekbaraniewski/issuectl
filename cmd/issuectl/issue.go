package cli

import (
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initAddRepoToIssueCommand(rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "addRepo [repo name] [issueID]",
		Short: "Add repo to issue. Clones repository to issue workdir and sets up branch.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoName := args[0]
			issueID := args[1] // TODO: this should be also figured out from context to just run `i addRepo someOtherRepo` while inside issue workdir

			return issuectl.AddRepoToIssue(repoName, issuectl.IssueID(issueID))

		},
	}

	rootCmd.AddCommand(listCmd)
}
