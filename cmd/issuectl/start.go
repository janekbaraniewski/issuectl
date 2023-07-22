package cli

import (
	"errors"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

var Flags issuectl.CLIOverwrites

func initStartCommand(rootCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:                "start [issue number]",
		Short:              "Start work on issue",
		Long:               `Create issue work directory. Clone all repositories from current profile. Create branches.`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly 1 argument - issue id")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			if err := issuectl.StartWorkingOnIssue(config, &Flags, issuectl.IssueID(args[0])); err != nil {
				issuectl.Log.Infof("Error!! -> %v", err)
				return err
			}

			return nil
		},
	}

	startCmd.PersistentFlags().StringSliceVarP(
		&Flags.Repos,
		"repos",
		"r",
		[]string{},
		"A list of repositories to clone",
	)

	startCmd.PersistentFlags().StringVarP(
		&Flags.Profile,
		"profile",
		"p",
		"",
		"Name of profile to use for command",
	)

	startCmd.PersistentFlags().StringVarP(
		&Flags.Backend,
		"backend",
		"b",
		"",
		"Name of issue backend to use for command",
	)

	rootCmd.AddCommand(startCmd)
}
