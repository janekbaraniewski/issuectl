package cli

import (
	"errors"
	"fmt"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

type CLIOverwrites struct {
	Repos   []string
	Profile string
	Backend string
}

var Flags CLIOverwrites

func MergeConfigWithOverwrites(conf issuectl.IssuectlConfig, overwrites *CLIOverwrites) (issuectl.IssuectlConfig, error) {
	conf = conf.GetInMemory()

	if overwrites.Profile != "" {
		if err := conf.UseProfile(issuectl.ProfileName(overwrites.Profile)); err != nil {
			return conf, err
		}
	}

	overwriteProfile := conf.GetProfile(conf.GetCurrentProfile())
	if overwriteProfile == nil {
		return conf, fmt.Errorf("Failed - profile %v not defined.", overwrites.Profile)
	}
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
	return conf, nil
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
			config, err := MergeConfigWithOverwrites(issuectl.LoadConfig(), &Flags)
			if err != nil {
				return err
			}
			if err := issuectl.StartWorkingOnIssue(config, issuectl.IssueID(args[0])); err != nil {
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
