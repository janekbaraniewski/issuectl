package cli

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func askForGitUser() (issuectl.GitUser, error) {
	answers := struct {
		Name   string
		Email  string
		SSHKey string
	}{}
	prompt := []*survey.Question{
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
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.GitUser{}, err
	}
	return issuectl.GitUser{
		Name:   issuectl.GitUserName(answers.Name),
		Email:  answers.Email,
		SSHKey: answers.SSHKey,
	}, nil
}

func askForBackend() (issuectl.BackendConfig, error) {
	answers := struct {
		Type  string
		Token string
	}{}
	prompt := []*survey.Question{
		{
			Name: "Type",
			Prompt: &survey.Select{
				Message: "Select backend type:",
				Options: []string{string(issuectl.BackendGithub), string(issuectl.BackendGitLab)},
			},
			Validate: survey.Required,
		},
		{
			Name: "Token",
			Prompt: &survey.Password{
				Message: "Enter backend token:",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.BackendConfig{}, err
	}
	return issuectl.BackendConfig{
		Name:  issuectl.BackendConfigName("default"),
		Type:  issuectl.BackendType(answers.Type),
		Token: base64.RawStdEncoding.EncodeToString([]byte(answers.Token)),
	}, nil
}

func askForProfile() (issuectl.Profile, error) {
	answers := struct {
		Workdir string
	}{}
	prompt := []*survey.Question{
		{
			Name: "workdir",
			Prompt: &survey.Input{
				Message: "Enter working directory for profile:",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.Profile{}, err
	}
	return issuectl.Profile{
		Name:    issuectl.ProfileName("default"),
		WorkDir: answers.Workdir,
	}, nil
}

func askForRepo() (issuectl.RepoConfig, error) {
	answers := struct {
		Name    string
		Owner   string
		RepoURL string
	}{}
	prompt := []*survey.Question{
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
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.RepoConfig{}, err
	}
	return issuectl.RepoConfig{
		Name:    issuectl.RepoConfigName(answers.Name),
		Owner:   answers.Owner,
		RepoURL: issuectl.RepoURL(answers.RepoURL),
	}, nil
}

func initInitConfigCommand(rootCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := issuectl.DefaultConfigFilePath
			_, err := os.Stat(configPath)
			if err == nil {
				return fmt.Errorf("config file already exists at %s", configPath)
			}

			gitUser, err := askForGitUser()
			if err != nil {
				return err
			}

			var backend issuectl.BackendConfig
			configureBackend := false
			prompt := &survey.Confirm{
				Message: "Do you want to configure a backend?",
			}
			survey.AskOne(prompt, &configureBackend)
			if configureBackend {
				backend, err = askForBackend()
				if err != nil {
					return err
				}
			}

			profile, err := askForProfile()
			if err != nil {
				return err
			}

			repo, err := askForRepo()
			if err != nil {
				return err
			}

			profile.GitUserName = gitUser.Name
			profile.AddRepository(&repo)
			profile.Backend = backend.Name
			profile.DefaultRepository = repo.Name

			config := issuectl.IssuectlConfig{
				Repositories: map[issuectl.RepoConfigName]issuectl.RepoConfig{
					repo.Name: repo,
				},
				Backends: map[issuectl.BackendConfigName]issuectl.BackendConfig{
					backend.Name: backend,
				},
				GitUsers: map[issuectl.GitUserName]issuectl.GitUser{
					gitUser.Name: gitUser,
				},
				Profiles: map[issuectl.ProfileName]issuectl.Profile{
					profile.Name: profile,
				},
			}

			return config.Save()
		},
	}
	rootCmd.AddCommand(initCmd)
}
