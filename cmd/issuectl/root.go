package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issuectl",
		Version: version,
		Short:   "issuectl",
		Long: `issuectl
	issuectl
		issuectl`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	initStartCommand(cmd)
	initFinishCommand(cmd)
	initRepositoriesCommand(cmd)
	initProfileCommand(cmd)
	initOpenPullRequestCommand(cmd)
	initConfigCommand(cmd)
	return cmd
}

func Execute(version string) {
	if err := RootCmd(version).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
