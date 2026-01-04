package cmd

import (
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	Long:  `Apply database migrations using Atlas.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateDatabase()
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Long:  `Seed the database with test data using Ent.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.SeedDatabase()
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup database",
	Long:  `Create a database backup using mysqldump.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.BackupDatabase()
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore [backup-file]",
	Short: "Restore database from backup",
	Long:  `Restore database from a backup SQL file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.RestoreDatabase(args[0])
	},
}

var dbDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy database (init, install, start, migrate, seed)",
	Long:  `Full database deployment workflow.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.DeployDatabase()
	},
}

var dbRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart database service",
	Long:  `Restart the MySQL database service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.RestartDatabase()
	},
}

var dbStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop database service",
	Long:  `Stop the MySQL database service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.StopDatabase()
	},
}

var dbUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall database service",
	Long:  `Uninstall the MySQL database service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.UninstallDatabase()
	},
}

func init() {
	databaseCmd.AddCommand(migrateCmd)
	databaseCmd.AddCommand(seedCmd)
	databaseCmd.AddCommand(backupCmd)
	databaseCmd.AddCommand(restoreCmd)
	databaseCmd.AddCommand(dbDeployCmd)
	databaseCmd.AddCommand(dbRestartCmd)
	databaseCmd.AddCommand(dbStopCmd)
	databaseCmd.AddCommand(dbUninstallCmd)
}
