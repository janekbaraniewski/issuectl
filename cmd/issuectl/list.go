package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"

	"github.com/spf13/cobra"
)

func initListIssuesCommand(rootCmd *cobra.Command) {
	var listIssuesCmd = &cobra.Command{
		Use:   "list",
		Short: "List all issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			issues := issuectl.LoadConfig().GetIssues() // Load your configuration here

			if len(issues) == 0 {
				fmt.Println("No issues found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "ID\tName\tRepository Backend\tIssue Backend\tRepositories\tProfile\t")
			for issueID, issue := range issues {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t\n", issueID, issue.Name, issue.RepoBackend, issue.IssueBackend, issue.Repositories, issue.Profile)
			}
			w.Flush()

			return nil
		},
	}

	rootCmd.AddCommand(listIssuesCmd)
}
