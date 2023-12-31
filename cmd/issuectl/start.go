package cli

import (
	"errors"
	"fmt"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

type CLIOverwrites struct {
	Repos        []string
	Profile      string
	IssueBackend string
	RepoBackend  string
	IssueName    string
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
	if overwrites.IssueBackend != "" {
		overwriteProfile.IssueBackend = issuectl.BackendConfigName(overwrites.IssueBackend)
	}
	if overwrites.Repos != nil {
		for _, repoName := range overwrites.Repos {
			overwriteProfile.Repositories = append(overwriteProfile.Repositories, issuectl.RepoConfigName(repoName))
		}
	}

	if err := conf.UpdateProfile(overwriteProfile); err != nil {
		return nil, err
	}
	return conf, nil
}

func initStartCommand(rootCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:                "start [issueID]",
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
			if err := issuectl.StartWorkingOnIssue(Flags.IssueName, config.GetPersistent(), issuectl.IssueID(args[0])); err != nil {
				// TODO: rollback changes made by StartWorkingOnIssue
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
		&Flags.IssueBackend,
		"issue-backend",
		"",
		"",
		"Name of issue backend to use for command",
	)

	startCmd.PersistentFlags().StringVarP(
		&Flags.RepoBackend,
		"repo-backend",
		"",
		"",
		"Name of repo backend to use for command",
	)

	startCmd.PersistentFlags().StringVarP(
		&Flags.IssueName,
		"name",
		"n",
		"",
		"Custom issue name to use [defaults to IssueID]. IssueID will be added as prefix of custom name.",
	)

	rootCmd.AddCommand(startCmd)
}
