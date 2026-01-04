package cmd

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Service management commands",
	Long:  `Manage API service lifecycle (install, start, stop, restart, uninstall).`,
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
