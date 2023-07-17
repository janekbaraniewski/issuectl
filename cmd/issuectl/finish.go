package cli

import (
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initFinishCommand(rootCmd *cobra.Command) {
	finishCmd := &cobra.Command{
		Use:                "finish",
		Short:              "finish",
		Long:               "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := issuectl.FinishWorkingOnIssue("test-issue-id"); err != nil {
				issuectl.Log.Infof("Error!! -> %v", err)
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(finishCmd)
}
