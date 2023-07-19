package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	ShortDescription string = "issuectl helps you manage separate environments for work on multiple issues"
	LongDescription  string = `
issuectl helps managing separate environments for work on multiple issues.

Start work on issue:
	issuectl start [issue_number]

Open PR and link it to issue
	issuectl openpr [issue_number]

Finish work, close the issue
	issuectl finish [issue_number]
`
)

func RootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issuectl",
		Version: version,
		Short:   ShortDescription,
		Long:    LongDescription,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	initStartCommand(cmd)
	initFinishCommand(cmd)
	initOpenPullRequestCommand(cmd)
	initConfigCommand(cmd)
	initInitConfigCommand(cmd)
	initListIssuesCommand(cmd)
	initOpenIssueCommand(cmd)
	return cmd
}

func Execute(version string) {
	if err := RootCmd(version).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
