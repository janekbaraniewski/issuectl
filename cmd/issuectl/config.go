package cli

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func initConfigCommand(rootCmd *cobra.Command) {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage config",
		Long:  `Manage issuectl config.`,
	}

	initPrintConfigCommand(configCmd)
	initBackendCommand(configCmd)
	initRepositoriesCommand(configCmd)
	initProfileCommand(configCmd)
	initGitUserCommand(configCmd)
	rootCmd.AddCommand(configCmd)
}

func initPrintConfigCommand(root *cobra.Command) {
	getConfigCmd := &cobra.Command{
		Use:   "get",
		Short: "Get config",
		Long:  `Prints full currently sellected config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			confYaml, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			fmt.Printf("%v", string(confYaml))
			return nil
		},
	}
	root.AddCommand(getConfigCmd)
}

func initBackendCommand(rootCmd *cobra.Command) {
	backendCmd := &cobra.Command{
		Use:   "backend",
		Short: "Manage backends",
		Long:  `Manage issuectl backends. You can list, add, delete, and use backends.`,
	}

	initBackendListCommand(backendCmd)
	initBackendAddCommand(backendCmd)
	initBackendDeleteCommand(backendCmd)
	initBackendUseCommand(backendCmd)

	rootCmd.AddCommand(backendCmd)
}

func initBackendListCommand(rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all backends",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "NAME\tTYPE\t")
			for _, backend := range issuectl.LoadConfig().GetBackends() {
				fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t", backend.Name, backend.Type))
			}
			w.Flush()
		},
	}

	rootCmd.AddCommand(listCmd)
}

func initBackendAddCommand(rootCmd *cobra.Command) {
	type _flags struct {
		GitHubApi    string
		GitHubToken  string
		GitLabApi    string
		GitLabToken  string
		JiraHost     string
		JiraToken    string
		JiraUsername string
	}

	var flags *_flags = &_flags{}

	addCmd := &cobra.Command{
		Use:   "add [name] [config]",
		Short: "Add a new backend",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig().GetPersistent()
			backendName := args[0]
			_backendType := args[1]

			backendType := issuectl.BackendType(_backendType)

			newBackend := issuectl.BackendConfig{
				Name: issuectl.BackendConfigName(backendName),
				Type: backendType,
			}

			switch backendType {

			case issuectl.BackendGithub:
				token := base64.RawStdEncoding.EncodeToString([]byte(flags.GitHubToken))
				githubConfig := &issuectl.GitHubConfig{
					Host:  flags.GitHubApi,
					Token: token,
				}
				newBackend.GitHub = githubConfig

			case issuectl.BackendGitLab:
				token := base64.RawStdEncoding.EncodeToString([]byte(flags.GitLabToken))
				gitlabConfig := &issuectl.GitLabConfig{
					Host:  flags.GitLabApi,
					Token: token,
				}
				newBackend.GitLab = gitlabConfig

			case issuectl.BackendJira:
				token := base64.RawStdEncoding.EncodeToString([]byte(flags.JiraToken))
				jiraBackend := &issuectl.JiraConfig{
					Host:     flags.JiraHost,
					Token:    token,
					Username: flags.JiraUsername,
				}
				newBackend.Jira = jiraBackend
			}
			return config.AddBackend(&newBackend)
		},
	}

	addCmd.PersistentFlags().StringVarP(
		&flags.GitHubApi,
		"github-api",
		"",
		"https://api.github.com/",
		"GitHub API URL",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.GitHubToken,
		"github-token",
		"",
		"",
		"GitHub API Auth Token",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.GitLabApi,
		"gitlab-api",
		"",
		"",
		"GitLab API URL",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.GitLabToken,
		"gitlab-token",
		"",
		"",
		"GitLab API Token",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.JiraHost,
		"jira-host",
		"",
		"",
		"Jira API URL",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.JiraToken,
		"jira-token",
		"",
		"",
		"Jira API Token",
	)

	addCmd.PersistentFlags().StringVarP(
		&flags.JiraUsername,
		"jira-username",
		"",
		"",
		"Jira API Username",
	)

	rootCmd.AddCommand(addCmd)
}

func initBackendDeleteCommand(rootCmd *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a backend",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig().GetPersistent()
			backendName := args[0]
			return config.DeleteBackend(issuectl.BackendConfigName(backendName))
		},
	}

	rootCmd.AddCommand(deleteCmd)
}

func initBackendUseCommand(rootCmd *cobra.Command) {
	useCmd := &cobra.Command{
		Use:   "use [name]",
		Short: "Use a backend",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// config := issuectl.LoadConfig()
		},
	}

	rootCmd.AddCommand(useCmd)
}

func initRepoListCommand(rootCmd *cobra.Command) {
	repoListCmd := &cobra.Command{
		Use:                "list",
		Short:              "List all repositories",
		Long:               `List all repositories`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "NAME\tOWNER\tURL\t")
			for _, repo := range issuectl.LoadConfig().GetRepositories() {
				fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t", repo.Name, repo.Owner, repo.RepoURL))
			}
			w.Flush()
		},
	}

	rootCmd.AddCommand(repoListCmd)
}

func initRepoAddCommand(rootCmd *cobra.Command) {
	repoAddCmd := &cobra.Command{
		Use:                "add",
		Short:              "Add a new repository",
		Long:               `Add a new repository`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 || len(args) > 3 {
				return errors.New("you must provide exactly 3 arguments - owner, name and url of repository")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := issuectl.LoadConfig().GetPersistent()
			repoConfig := &issuectl.RepoConfig{
				Owner:   args[0],
				Name:    issuectl.RepoConfigName(args[1]),
				RepoURL: issuectl.RepoURL(args[2]),
			}
			return conf.AddRepository(repoConfig)
		},
	}

	rootCmd.AddCommand(repoAddCmd)
}

func initRepositoriesCommand(rootCmd *cobra.Command) {
	repoCmd := &cobra.Command{
		Use:                "repository",
		Aliases:            []string{"repo"},
		Short:              "Manage repositories",
		Long:               `Manage repositories`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	initRepoListCommand(repoCmd)
	initRepoAddCommand(repoCmd)
	rootCmd.AddCommand(repoCmd)
}

func initGitUserListCommand(rootCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all Git users",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "NAME\tEMAIL\tSSH KEY\t")
			for _, user := range issuectl.LoadConfig().GetGitUsers() {
				fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t", user.Name, user.Email, user.SSHKey))
			}
			w.Flush()
		},
	}

	rootCmd.AddCommand(listCmd)
}

func initGitUserAddCommand(rootCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [username] [email] [sshKey]",
		Short: "Add a new Git user",
		Args:  cobra.ExactArgs(3), // Expects exactly 3 arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := issuectl.LoadConfig().GetPersistent()
			gitUser := &issuectl.GitUser{
				Name:   issuectl.GitUserName(args[0]),
				Email:  args[1],
				SSHKey: args[2],
			}
			return conf.AddGitUser(gitUser)
		},
	}

	rootCmd.AddCommand(addCmd)
}

func initGitUserDeleteCommand(rootCmd *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Git user",
		Args:  cobra.ExactArgs(1), // Expects exactly 1 argument
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := issuectl.LoadConfig()
			return conf.DeleteGitUser(issuectl.GitUserName(args[0]))
		},
	}

	rootCmd.AddCommand(deleteCmd)
}

func initGitUserCommand(rootCmd *cobra.Command) {
	userCmd := &cobra.Command{
		Use:   "gituser",
		Short: "Manage Git users",
		Long:  `Manage Git users. You can list, add, delete Git users.`,
	}

	initGitUserListCommand(userCmd)
	initGitUserAddCommand(userCmd)
	initGitUserDeleteCommand(userCmd)

	rootCmd.AddCommand(userCmd)
}
