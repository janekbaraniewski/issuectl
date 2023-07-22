package cli

import (
	"errors"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

var Flags issuectl.CLIOverwrites

func MergeConfigWithOverwrites(conf issuectl.IssuectlConfig, overwrites *issuectl.CLIOverwrites) issuectl.IssuectlConfig {
	conf = conf.GetInMemory()

	if overwrites.Profile != "" {
		conf.UseProfile(issuectl.ProfileName(overwrites.Profile))
	}

	overwriteProfile := conf.GetProfile(conf.GetCurrentProfile())
	if overwrites.Backend != "" {
		overwriteProfile.Backend = issuectl.BackendConfigName(overwrites.Profile)
	}
	if overwrites.Repos != nil {
		repos := []issuectl.RepoConfigName{}

		for _, repoName := range overwrites.Repos {
			repos = append(repos, issuectl.RepoConfigName(repoName))
		}

		overwriteProfile.Repositories = repos
	}

	conf.UpdateProfile(overwriteProfile)
	return conf
}

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
			config = MergeConfigWithOverwrites(config, &Flags)
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
