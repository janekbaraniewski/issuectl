package cli

import (
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initRepoListCommand(rootCmd *cobra.Command) {
	repoListCmd := &cobra.Command{
		Use:                "list",
		Short:              "",
		Long:               "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return issuectl.ListRepositories()
		},
	}

	rootCmd.AddCommand(repoListCmd)
}

func initRepoAddCommand(rootCmd *cobra.Command) {
	repoAddCmd := &cobra.Command{
		Use:                "add",
		Short:              "",
		Long:               "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return issuectl.AddRepository()
		},
	}

	rootCmd.AddCommand(repoAddCmd)
}

func initRepositoriesCommand(rootCmd *cobra.Command) {
	repoCmd := &cobra.Command{
		Use:                "repository",
		Aliases:            []string{"repo"},
		Short:              "",
		Long:               "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	initRepoListCommand(repoCmd)
	initRepoAddCommand(repoCmd)
	rootCmd.AddCommand(repoCmd)
}
