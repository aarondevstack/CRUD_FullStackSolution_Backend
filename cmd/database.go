package cmd

import (
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database management commands",
	Long:  `Manage MySQL database lifecycle, migrations, seeding, backup, and restore.`,
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
