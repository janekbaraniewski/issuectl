package cli

import (
	"fmt"
	"os"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func initInitConfigCommand(rootCmd *cobra.Command) {
	configCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := issuectl.DefaultConfigFilePath // Specify your config path here
			_, err := os.Stat(configPath)
			if err == nil {
				return fmt.Errorf("config file already exists at %s", configPath)
			}

			var gitUser issuectl.GitUser
			var backend issuectl.BackendConfig
			var profile issuectl.Profile
			var repo issuectl.RepoConfig

			gitUserSurvey := []*survey.Question{
				{
					Name: "Name",
					Prompt: &survey.Input{
						Message: "Enter Git user name:",
					},
					Validate: survey.Required,
				},
				{
					Name: "Email",
					Prompt: &survey.Input{
						Message: "Enter Git user email:",
					},
					Validate: survey.Required,
				},
				{
					Name: "SSHKey",
					Prompt: &survey.Input{
						Message: "Enter SSH key path:",
					},
					Validate: survey.Required,
				},
			}

			// Define survey questions for BackendConfig
			backendSurvey := []*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Enter backend config name:",
					},
					Validate: survey.Required,
				},
				{
					Name: "type",
					Prompt: &survey.Select{
						Message: "Select backend type:",
						Options: []string{string(issuectl.BackendGithub), string(issuectl.BackendGitLab)},
					},
					Validate: survey.Required,
				},
			}

			// Define survey questions for Profile
			profileSurvey := []*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Enter profile name:",
					},
					Validate: survey.Required,
				},
				{
					Name: "workdir",
					Prompt: &survey.Input{
						Message: "Enter working directory for profile:",
					},
					Validate: survey.Required,
				},
			}

			// Define survey questions for RepoConfig
			repoSurvey := []*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Enter repository name:",
					},
					Validate: survey.Required,
				},
				{
					Name: "owner",
					Prompt: &survey.Input{
						Message: "Enter repository owner:",
					},
					Validate: survey.Required,
				},
				{
					Name: "repourl",
					Prompt: &survey.Input{
						Message: "Enter repository URL:",
					},
					Validate: survey.Required,
				},
			}

			// Perform the surveys
			survey.Ask(gitUserSurvey, &gitUser)
			survey.Ask(backendSurvey, &backend)
			survey.Ask(profileSurvey, &profile)
			survey.Ask(repoSurvey, &repo)

			fmt.Printf("%v", gitUser)

			// Create and save IssuectlConfig
			config := issuectl.IssuectlConfig{
				Repositories: map[issuectl.RepoConfigName]issuectl.RepoConfig{
					issuectl.RepoConfigName(repo.Name): repo,
				},
				Backends: map[issuectl.BackendConfigName]issuectl.BackendConfig{
					issuectl.BackendConfigName(backend.Name): backend,
				},
				GitUsers: map[issuectl.GitUserName]issuectl.GitUser{
					issuectl.GitUserName(gitUser.Name): gitUser,
				},
				Profiles: map[issuectl.ProfileName]issuectl.Profile{
					issuectl.ProfileName(profile.Name): profile,
				},
			}
			fmt.Printf("Saving\n")
			// Save the config here
			err = config.Save()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			return err
		},
	}

	rootCmd.AddCommand(configCmd)
}
