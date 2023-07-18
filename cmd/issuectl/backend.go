package cli

import (
	"fmt"

	issuectl "github.com/janekbaraniewski/issuectl/pkg"
	"github.com/spf13/cobra"
)

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
			config := issuectl.LoadConfig()
			for _, backend := range config.Backends {
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
		Run: func(cmd *cobra.Command, args []string) {
			config := issuectl.LoadConfig()
			backendName := args[0]
			backendType := args[1]
			newBackend := issuectl.BackendConfig{
				Name: issuectl.BackendConfigName(backendName),
				Type: issuectl.BackendType(backendType),
			}
			config.Backends = append(config.Backends, newBackend)
			if err := config.Save(); err != nil {
				issuectl.Log.Infof("Failed to save config: %v", err)
				return
			}
		},
	}

	rootCmd.AddCommand(addCmd)
}

func initBackendDeleteCommand(rootCmd *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a backend",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config := issuectl.LoadConfig()
			backendName := args[0]
			for i, backend := range config.Backends {
				if string(backend.Name) == backendName {
					config.Backends = append(config.Backends[:i], config.Backends[i+1:]...)
					break
				}
			}
			if err := config.Save(); err != nil {
				issuectl.Log.Infof("Failed to save config: %v", err)
				return
			}
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
			config := issuectl.LoadConfig()
			backendName := args[0]
			for _, backend := range config.Backends {
				if string(backend.Name) == backendName {
					// TODO: Set the current backend to the specified backend
					// This will depend on how you're managing the current backend in your application
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