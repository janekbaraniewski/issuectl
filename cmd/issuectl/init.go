package cli

import (
	"fmt"
	"os"
	"strconv"

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
		Type string
	}{}
	prompt := []*survey.Question{
		{
			Name: "Type",
			Prompt: &survey.Select{
				Message: "Select backend type:",
				Options: []string{string(issuectl.BackendGithub), string(issuectl.BackendGitLab), string(issuectl.BackendJira)},
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.BackendConfig{}, err
	}
	switch answers.Type {

	case string(issuectl.BackendGithub):
		ghConf, err := askForGitHubConfig()
		if err != nil {
			return issuectl.BackendConfig{}, err
		}
		return issuectl.BackendConfig{
			Name:   issuectl.BackendConfigName("default"),
			Type:   issuectl.BackendType(answers.Type),
			GitHub: &ghConf,
		}, nil

	case string(issuectl.BackendGitLab):
		glConf, err := askForGitLabConfig()
		if err != nil {
			return issuectl.BackendConfig{}, err
		}
		return issuectl.BackendConfig{
			Name:   issuectl.BackendConfigName("default"),
			Type:   issuectl.BackendType(answers.Type),
			GitLab: &glConf,
		}, nil

	case string(issuectl.BackendJira):
		jiraConf, err := askForJiraConfig()
		if err != nil {
			return issuectl.BackendConfig{}, err
		}
		return issuectl.BackendConfig{
			Name: issuectl.BackendConfigName("default"),
			Type: issuectl.BackendType(answers.Type),
			Jira: &jiraConf,
		}, nil

	default:
		return issuectl.BackendConfig{
			Name: issuectl.BackendConfigName("default"),
			Type: issuectl.BackendType(answers.Type),
			// Token: base64.RawStdEncoding.EncodeToString([]byte(answers.Token)),
		}, nil

	}
}

func askForGitHubConfig() (issuectl.GitHubConfig, error) {
	answers := struct {
		Host     string
		Token    string
		Username string
	}{}
	prompt := []*survey.Question{
		{
			Name:   "Host",
			Prompt: &survey.Input{Message: "Enter GitHub Host (Skip for https://api.github.com/):"},
		},
		{
			Name:     "Token",
			Prompt:   &survey.Password{Message: "Enter GitHub Token:"},
			Validate: survey.Required,
		},
		{
			Name:     "Username",
			Prompt:   &survey.Input{Message: "Enter GitHub Username:"},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.GitHubConfig{}, err
	}
	host := answers.Host
	if host == "" {
		host = "https://api.github.com/"
	}
	return issuectl.GitHubConfig{
		Host:     answers.Host,
		Token:    answers.Token,
		Username: answers.Username,
	}, nil
}

func askForGitLabConfig() (issuectl.GitLabConfig, error) {
	answers := struct {
		Host   string
		Token  string
		UserID string
	}{}
	prompt := []*survey.Question{
		{
			Name:     "Host",
			Prompt:   &survey.Input{Message: "Enter GitLab Host:"},
			Validate: survey.Required,
		},
		{
			Name:     "Token",
			Prompt:   &survey.Password{Message: "Enter GitLab Token:"},
			Validate: survey.Required,
		},
		{
			Name:     "UserID",
			Prompt:   &survey.Input{Message: "Enter GitLab UserID:"},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.GitLabConfig{}, err
	}
	userID, err := strconv.Atoi(answers.UserID)
	if err != nil {
		return issuectl.GitLabConfig{}, err
	}
	return issuectl.GitLabConfig{
		Host:   answers.Host,
		Token:  answers.Token,
		UserID: userID,
	}, nil
}

func askForJiraConfig() (issuectl.JiraConfig, error) {
	answers := struct {
		Host     string
		Token    string
		Username string
	}{}
	prompt := []*survey.Question{
		{
			Name:     "Host",
			Prompt:   &survey.Input{Message: "Enter Jira Host:"},
			Validate: survey.Required,
		},
		{
			Name:     "Token",
			Prompt:   &survey.Password{Message: "Enter Jira Token:"},
			Validate: survey.Required,
		},
		{
			Name:     "Username",
			Prompt:   &survey.Input{Message: "Enter Jira Username:"},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(prompt, &answers); err != nil {
		return issuectl.JiraConfig{}, err
	}
	return issuectl.JiraConfig{
		Host:     answers.Host,
		Token:    answers.Token,
		Username: answers.Username,
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
			if err := survey.AskOne(prompt, &configureBackend); err != nil {
				return err
			}
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

			if err := profile.AddRepository(repo.Name); err != nil {
				return err
			}

			profile.IssueBackend = backend.Name
			profile.RepoBackend = backend.Name
			profile.DefaultRepository = repo.Name

			config := issuectl.GetConfig(
				profile.Name,
				map[issuectl.RepoConfigName]*issuectl.RepoConfig{
					repo.Name: &repo,
				},
				map[issuectl.BackendConfigName]*issuectl.BackendConfig{
					backend.Name: &backend,
				},
				map[issuectl.GitUserName]*issuectl.GitUser{
					gitUser.Name: &gitUser,
				},
				map[issuectl.ProfileName]*issuectl.Profile{
					profile.Name: &profile,
				},
			)

			return config.GetPersistent().Save()
		},
	}
	rootCmd.AddCommand(initCmd)
}
