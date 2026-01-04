package cmd

import (
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API service management",
	Long:  `Manage API service lifecycle: serve, install, start, stop, restart, status, uninstall, deploy.`,
}

func init() {
	servicesCmd.AddCommand(apiCmd)
}
