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

			gitUserSurvey := []*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Enter Git user name:",
					},
					Validate: survey.Required,
				},
				{
					Name: "email",
					Prompt: &survey.Input{
						Message: "Enter Git user email:",
					},
					Validate: survey.Required,
				},
				{
					Name: "sshKey",
					Prompt: &survey.Input{
						Message: "Enter SSH key path:",
					},
					Validate: survey.Required,
				},
			}

			// Define survey questions for BackendConfig
			backendSurvey := []*survey.Question{
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
			gitUserAnswers := struct {
				Name   string
				Email  string
				SSHKey string
			}{}
			survey.Ask(gitUserSurvey, &gitUserAnswers)
			backendAnswers := struct {
				Type string
			}{}
			survey.Ask(backendSurvey, &backendAnswers)
			profileAnswers := struct {
				Workdir string
			}{}
			survey.Ask(profileSurvey, &profileAnswers)
			repoAnswers := struct {
				Name    string
				Owner   string
				RepoURL string
			}{}
			survey.Ask(repoSurvey, &repoAnswers)

			// Create and save IssuectlConfig
			config := issuectl.IssuectlConfig{
				Repositories: map[issuectl.RepoConfigName]issuectl.RepoConfig{
					issuectl.RepoConfigName(repoAnswers.Name): {
						Name:    issuectl.RepoConfigName(repoAnswers.Name),
						Owner:   repoAnswers.Owner,
						RepoURL: issuectl.RepoURL(repoAnswers.RepoURL),
					},
				},
				Backends: map[issuectl.BackendConfigName]issuectl.BackendConfig{
					issuectl.BackendConfigName("default"): {
						Name: issuectl.BackendConfigName("default"),
						Type: issuectl.BackendType(backendAnswers.Type),
					},
				},
				GitUsers: map[issuectl.GitUserName]issuectl.GitUser{
					issuectl.GitUserName(gitUserAnswers.Name): {
						Name:   issuectl.GitUserName(gitUserAnswers.Name),
						Email:  gitUserAnswers.Email,
						SSHKey: gitUserAnswers.SSHKey,
					},
				},
				Profiles: map[issuectl.ProfileName]issuectl.Profile{
					issuectl.ProfileName("default"): {
						Name:    issuectl.ProfileName("default"),
						WorkDir: profileAnswers.Workdir,
					},
				},
			}

			return config.Save()
		},
	}

	rootCmd.AddCommand(configCmd)
}
