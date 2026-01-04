package database

import (
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/mysql"
)

// InitMySQL initializes MySQL database
func InitMySQL() error {
	return mysql.InitMySQL()
}

// InstallMySQL installs MySQL as a system service
func InstallMySQL() error {
	return mysql.InstallMySQL()
}

// UninstallMySQL uninstalls MySQL system service
func UninstallMySQL() error {
	return mysql.UninstallMySQL()
}

// StartMySQL starts the MySQL service
func StartMySQL() error {
	return mysql.StartMySQL()
}

// StopMySQL stops the MySQL service
func StopMySQL() error {
	return mysql.StopMySQL()
}

// StatusMySQL checks MySQL service status
func StatusMySQL() error {
	return mysql.StatusMySQL()
}

// RestartMySQL restarts the MySQL service
func RestartMySQL() error {
	return mysql.RestartMySQL()
}
