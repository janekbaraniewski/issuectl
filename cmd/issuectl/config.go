package cli

import (
	"fmt"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

func initConfigCommand(rootCmd *cobra.Command) {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage config",
		Long:  `Manage issuectl config.`,
	}
	initPrintConfigCommand(configCmd)
	initBackendCommand(configCmd)
	rootCmd.AddCommand(configCmd)
}

func initPrintConfigCommand(root *cobra.Command) {
	getConfigCmd := &cobra.Command{
		Use:   "get",
		Short: "Get config",
		Long:  `Prints full currently sellected config`,
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
