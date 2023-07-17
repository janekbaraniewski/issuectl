package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issuectl",
		Short: "issuectl",
		Long: `issuectl
	issuectl
		issuectl`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	initStartCommand(cmd)
	initFinishCommand(cmd)
	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
