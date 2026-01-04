package cmd

import (
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database"
	"github.com/spf13/cobra"
)

var mysqlInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize MySQL database",
	Long:  `Decompress MySQL ZIP, initialize data directory, and create initial users.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.InitMySQL()
	},
}

var mysqlInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install MySQL as system service",
	Long:  `Install MySQL as a system service using kardianos/service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.InstallMySQL()
	},
}

var mysqlUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall MySQL system service",
	Long:  `Uninstall MySQL system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.UninstallMySQL()
	},
}

var mysqlStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start MySQL service",
	Long:  `Start the MySQL system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.StartMySQL()
	},
}

var mysqlStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop MySQL service",
	Long:  `Stop the MySQL system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.StopMySQL()
	},
}

var mysqlStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check MySQL service status",
	Long:  `Query the status of the MySQL system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.StatusMySQL()
	},
}

var mysqlRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart MySQL service",
	Long:  `Restart the MySQL system service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.RestartMySQL()
	},
}

func init() {
	mysqlCmd.AddCommand(mysqlInitCmd)
	mysqlCmd.AddCommand(mysqlInstallCmd)
	mysqlCmd.AddCommand(mysqlUninstallCmd)
	mysqlCmd.AddCommand(mysqlStartCmd)
	mysqlCmd.AddCommand(mysqlStopCmd)
	mysqlCmd.AddCommand(mysqlStatusCmd)
	mysqlCmd.AddCommand(mysqlRestartCmd)
}
