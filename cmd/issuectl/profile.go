package cli

import (
	"fmt"

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
			profiles := issuectl.LoadConfig().GetProfiles()
			for _, profile := range profiles {
				fmt.Println(profile.Name)
			}
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
			repos := []*issuectl.RepoConfig{}
			for _, repoName := range repoList {
				repo := config.GetRepository(issuectl.RepoConfigName(repoName))
				repos = append(repos, &repo)
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
		&repoList,
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
			config := issuectl.LoadConfig()
			repoName := args[0]
			issuectl.Log.Infof("repoName - %v", repoName)
			profile := config.GetProfile(config.GetCurrentProfile())
			issuectl.Log.Infof("profile - %v", profile)
			repo := config.GetRepository(issuectl.RepoConfigName(repoName))
			issuectl.Log.Infof("repo - %v", repo)
			profile.AddRepository(&repo)
			issuectl.Log.Infof("profile after add -> %v", profile)
			if err := config.UpdateProfile(&profile); err != nil {
				return err
			}
			issuectl.Log.Infof("config -> %v", config)
			return config.Save()
		},
	}

	rootCmd.AddCommand(addRepoCmd)
}
