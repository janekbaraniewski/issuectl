package cli

import (
	"errors"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initFinishCommand(rootCmd *cobra.Command) {
	finishCmd := &cobra.Command{
		Use:                "finish [issue number]",
		Short:              "Cleanup resources and close issue",
		Long:               `Removes issue work directory. Closes issue in backend.`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly 1 argument - issue id")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := issuectl.FinishWorkingOnIssue(issuectl.IssueID(args[0])); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(finishCmd)
}
