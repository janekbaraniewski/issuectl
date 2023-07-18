package cli

import (
	"errors"
	"fmt"

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
			for _, backend := range issuectl.LoadConfig().GetBackends() {
				fmt.Println(backend.Name)
			}
		},
	}

	rootCmd.AddCommand(listCmd)
}

func initBackendAddCommand(rootCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [name] [config]",
		Short: "Add a new backend",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			backendName := args[0]
			backendType := args[1]
			newBackend := issuectl.BackendConfig{
				Name: issuectl.BackendConfigName(backendName),
				Type: issuectl.BackendType(backendType),
			}
			return config.AddBackend(&newBackend)
		},
	}

	rootCmd.AddCommand(addCmd)
}

func initBackendDeleteCommand(rootCmd *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a backend",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
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
		RunE: func(cmd *cobra.Command, args []string) error {
			config := issuectl.LoadConfig()
			return config.ListRepositories()
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
			conf := issuectl.LoadConfig()
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