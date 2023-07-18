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

	rootCmd.AddCommand(profileCmd)
}

func initProfileListCommand(rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			config := issuectl.LoadConfig()
			for _, profile := range config.Profiles {
				fmt.Println(profile.Name)
			}
		},
	}

	rootCmd.AddCommand(listCmd)
}

func initProfileAddCommand(rootCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [name] [workdir] [defaultrepository]",
		Short: "Add a new profile",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			profileName := args[0]
			workDir := args[1]
			defaultRepository := args[2]
			newProfile := &issuectl.Profile{
				Name:       issuectl.ProfileName(profileName),
				WorkDir:    workDir,
				Repository: issuectl.RepoConfigName(defaultRepository),
				Backend:    "github",
			}
			return config.AddProfile(newProfile)
		},
	}

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
		Run: func(cmd *cobra.Command, args []string) {
			config := issuectl.LoadConfig()
			profileName := args[0]
			for _, profile := range config.Profiles {
				if string(profile.Name) == profileName {
					config.CurrentProfile = profile.Name
					break
				}
			}
			if err := config.Save(); err != nil {
				issuectl.Log.Infof("Failed to save config: %v", err)
				return
			}
		},
	}

	rootCmd.AddCommand(useCmd)
}
