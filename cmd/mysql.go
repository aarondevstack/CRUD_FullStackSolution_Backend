package cmd

import (
	"github.com/spf13/cobra"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "MySQL service management",
	Long:  `Manage MySQL service lifecycle: init, install, start, stop, restart, status, uninstall.`,
}

func init() {
	databaseCmd.AddCommand(mysqlCmd)
}
