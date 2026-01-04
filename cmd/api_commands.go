package cmd

import (
	"fmt"
	"os"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the API server in foreground",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService(""); err != nil {
			fmt.Printf("Error running service: %v\n", err)
			os.Exit(1)
		}
	},
}

// installCmd
var installApiCmd = &cobra.Command{
	Use:   "install",
	Short: "Install API as system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService("install"); err != nil {
			fmt.Printf("Failed to install service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Service installed successfully!")
	},
}

// uninstallCmd
var uninstallApiCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall API system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService("uninstall"); err != nil {
			fmt.Printf("Failed to uninstall service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Service uninstalled successfully!")
	},
}

// startCmd
var startApiCmd = &cobra.Command{
	Use:   "start",
	Short: "Start API system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService("start"); err != nil {
			fmt.Printf("Failed to start service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Service started successfully!")
	},
}

// stopCmd
var stopApiCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop API system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService("stop"); err != nil {
			fmt.Printf("Failed to stop service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Service stopped successfully!")
	},
}

// restartCmd
var restartApiCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart API system service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := api.ManageService("restart"); err != nil {
			fmt.Printf("Failed to restart service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Service restarted successfully!")
	},
}

// statusCmd
var statusApiCmd = &cobra.Command{
	Use:   "status",
	Short: "Check API service status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking service status...")
		svcConfig := &service.Config{
			Name: "CRUDSolution_WebService",
			Option: service.KeyValue{
				"UserService": true,
			},
		}
		prg := api.NewAPIService()
		s, err := service.New(prg, svcConfig)
		if err != nil {
			fmt.Printf("Status check failed: %v\n", err)
			return
		}

		status, err := s.Status()
		if err != nil {
			fmt.Printf("Failed to get status: %v\n", err)
			return
		}

		switch status {
		case service.StatusRunning:
			fmt.Println("Service is running ✅")
		case service.StatusStopped:
			fmt.Println("Service is stopped ⏹️")
		case service.StatusUnknown:
			fmt.Println("Service status is unknown ❓")
		default:
			fmt.Printf("Service status code: %d\n", status)
		}
	},
}

func init() {
	servicesCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(serveCmd)
	apiCmd.AddCommand(installApiCmd)
	apiCmd.AddCommand(uninstallApiCmd)
	apiCmd.AddCommand(startApiCmd)
	apiCmd.AddCommand(stopApiCmd)
	apiCmd.AddCommand(restartApiCmd)
	apiCmd.AddCommand(statusApiCmd)
}
