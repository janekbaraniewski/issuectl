package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initProfileCommand(rootCmd *cobra.Command) {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
		Long:  `Manage issuectl profiles. You can list, add, delete, and use profiles.`,
	}

	initProfileListCommand(profileCmd)
	initProfileAddCommand(profileCmd)
	initProfileDeleteCommand(profileCmd)
	initProfileUseCommand(profileCmd)
	initProfileAddRepoCommand(profileCmd)

	rootCmd.AddCommand(profileCmd)
}

func initProfileListCommand(rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			config := issuectl.LoadConfig()
			profiles := config.GetProfiles()
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "NAME\tWORK DIR\tGIT USER\tBACKEND\tREPOSITORIES\t")
			for _, profile := range profiles {
				repos := []string{}
				for _, repoName := range profile.Repositories {
					repos = append(repos, string(*repoName))
				}
				fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t", profile.Name, profile.WorkDir, profile.GitUserName, profile.Backend, repos))
			}
			w.Flush()
		},
	}

	rootCmd.AddCommand(listCmd)
}

func initProfileAddCommand(rootCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [name] [workdir]",
		Short: "Add a new profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			profileName := args[0]
			workDir := args[1]
			repos := []*issuectl.RepoConfigName{}
			for _, repoName := range Flags.Repos {
				repos = append(repos, (*issuectl.RepoConfigName)(&repoName))
			}
			newProfile := &issuectl.Profile{
				Name:         issuectl.ProfileName(profileName),
				WorkDir:      workDir,
				Repositories: repos,
				Backend:      "github",
			}
			return config.AddProfile(newProfile)
		},
	}

	addCmd.PersistentFlags().StringSliceVarP(
		&Flags.Repos,
		"repos",
		"r",
		[]string{},
		"A list of repositories to clone",
	)

	rootCmd.AddCommand(addCmd)
}

func initProfileDeleteCommand(rootCmd *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			profileName := args[0]
			return config.DeleteProfile(issuectl.ProfileName(profileName))
		},
	}

	rootCmd.AddCommand(deleteCmd)
}

func initProfileUseCommand(rootCmd *cobra.Command) {
	useCmd := &cobra.Command{
		Use:   "use [name]",
		Short: "Use a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			profileName := args[0]
			return config.UseProfile(issuectl.ProfileName(profileName))
		},
	}

	rootCmd.AddCommand(useCmd)
}

func initProfileAddRepoCommand(rootCmd *cobra.Command) {
	addRepoCmd := &cobra.Command{
		Use:   "addRepo [repo name]",
		Short: "Add a new repository to current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig().GetPersistent()
			repoName := args[0]
			profile := config.GetProfile(config.GetCurrentProfile())
			profile.AddRepository((*issuectl.RepoConfigName)(&repoName))
			if err := config.UpdateProfile(&profile); err != nil {
				return err
			}

			return config.Save()
		},
	}

	rootCmd.AddCommand(addRepoCmd)
}
