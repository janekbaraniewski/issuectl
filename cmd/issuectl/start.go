package cli

import (
	"errors"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

var repoList []string

func initStartCommand(rootCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:                "start",
		Short:              "start",
		Long:               "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly 1 argument - issue id")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := issuectl.StartWorkingOnIssue(repoList, issuectl.IssueID(args[0])); err != nil {
				issuectl.Log.Infof("Error!! -> %v", err)
				return err
			}

			return nil
		},
	}

	startCmd.PersistentFlags().StringSliceVarP(
		&repoList,
		"repos",
		"r",
		[]string{},
		"A list of repositories to clone",
	)

	rootCmd.AddCommand(startCmd)
}
