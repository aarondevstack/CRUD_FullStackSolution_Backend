package cmd

import (
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/cli/services"
	"github.com/spf13/cobra"
)

var apiServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server in foreground",
	Long:  `Start the API server in foreground mode (blocking, terminate with Ctrl+C).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.ServeAPI()
	},
}

var apiInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install API as system service",
	Long:  `Install API as a system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.InstallAPI()
	},
}

var apiUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall API system service",
	Long:  `Uninstall API system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.UninstallAPI()
	},
}

var apiStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start API service",
	Long:  `Start the API system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.StartAPI()
	},
}

var apiStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop API service",
	Long:  `Stop the API system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.StopAPI()
	},
}

var apiStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check API service status",
	Long:  `Query the status of the API system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.StatusAPI()
	},
}

var apiRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart API service",
	Long:  `Restart the API system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.RestartAPI()
	},
}

var apiDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy API service (install and start)",
	Long:  `Deploy API service by installing and starting it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return services.DeployAPI()
	},
}

func init() {
	apiCmd.AddCommand(apiServeCmd)
	apiCmd.AddCommand(apiInstallCmd)
	apiCmd.AddCommand(apiUninstallCmd)
	apiCmd.AddCommand(apiStartCmd)
	apiCmd.AddCommand(apiStopCmd)
	apiCmd.AddCommand(apiStatusCmd)
	apiCmd.AddCommand(apiRestartCmd)
	apiCmd.AddCommand(apiDeployCmd)
}
