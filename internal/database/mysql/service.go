package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
	"github.com/kardianos/service"
)

const (
	mysqlServiceName        = "com.github.aarondevstack.CRUDSolution_Database"
	mysqlServiceDisplayName = "CRUD Solution MySQL Database"
	mysqlServiceDescription = "MySQL database service for CRUD Solution"
)

// MySQLService implements service.Interface
type MySQLService struct {
	mysqldPath string
	configPath string
	mysqlDir   string
	cmd        *exec.Cmd
}

// Start implements service.Interface.Start
// Starts mysqld with the configuration file
func (m *MySQLService) Start(s service.Service) error {
	// Set environment for MySQL libraries
	env := os.Environ()
	if runtime.GOOS == "darwin" {
		env = append(env, fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(m.mysqlDir, "lib")))
	}

	// Start mysqld with config file
	m.cmd = exec.Command(m.mysqldPath, fmt.Sprintf("--defaults-file=%s", m.configPath))
	m.cmd.Env = env

	if err := m.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mysqld: %w", err)
	}

	return nil
}

// Stop implements service.Interface.Stop
// Gracefully stops mysqld using mysqladmin shutdown
func (m *MySQLService) Stop(s service.Service) error {
	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base directory: %w", err)
	}

	mysqladminPath := filepath.Join(baseDir, "mysql", "bin", "mysqladmin")

	// Set environment for MySQL libraries
	env := os.Environ()
	if runtime.GOOS == "darwin" {
		env = append(env, fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(m.mysqlDir, "lib")))
	}

	shutdownCmd := exec.Command(mysqladminPath, "-u", "root", "-paarondevstack@2026", "shutdown")
	shutdownCmd.Env = env

	if err := shutdownCmd.Run(); err != nil {
		// If graceful shutdown fails, try to kill the process
		if m.cmd != nil && m.cmd.Process != nil {
			m.cmd.Process.Kill()
		}
		return fmt.Errorf("failed to shutdown MySQL gracefully: %w", err)
	}

	return nil
}

// generateMySQLConfig generates mysql.ini configuration file
func generateMySQLConfig() error {
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base directory: %w", err)
	}

	mysqlDir := filepath.Join(baseDir, "mysql")
	dataDir := filepath.Join(mysqlDir, "data")
	logDir := filepath.Join(baseDir, "logs")
	pidFile := filepath.Join(mysqlDir, "mysqld.pid")

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	configContent := fmt.Sprintf(`[mysqld]
port=%d
basedir=%s
datadir=%s
pid-file=%s
log-error=%s/mysql-error.log
character-set-server=utf8mb4
default-storage-engine=INNODB
`,
		config.AppConfig.Database.Port,
		mysqlDir,
		dataDir,
		pidFile,
		logDir,
	)

	configPath := filepath.Join(mysqlDir, "mysql.ini")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write MySQL config: %w", err)
	}

	return nil
}

// createMySQLService creates a new service instance
func createMySQLService() (service.Service, error) {
	baseDir, err := config.GetBaseDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get base directory: %w", err)
	}

	mysqlDir := filepath.Join(baseDir, "mysql")
	mysqldPath := filepath.Join(mysqlDir, "bin", "mysqld")
	configPath := filepath.Join(mysqlDir, "mysql.ini")

	// Check if mysqld exists
	if _, err := os.Stat(mysqldPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mysqld not found. Run 'database mysql init' first")
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("MySQL config not found. Run 'database mysql install' first")
	}

	mysqlService := &MySQLService{
		mysqldPath: mysqldPath,
		configPath: configPath,
		mysqlDir:   mysqlDir,
	}

	svcConfig := &service.Config{
		Name:        mysqlServiceName,
		DisplayName: mysqlServiceDisplayName,
		Description: mysqlServiceDescription,
		Executable:  mysqldPath, // Direct path to mysqld
		Arguments:   []string{fmt.Sprintf("--defaults-file=%s", configPath)},
		Option: service.KeyValue{
			"UserService": true, // Use user mode for development (macOS/Linux)
		},
	}

	// Add environment variables for macOS
	if runtime.GOOS == "darwin" {
		svcConfig.Option["EnvironmentVariables"] = fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(mysqlDir, "lib"))
	}

	svc, err := service.New(mysqlService, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return svc, nil
}

// InstallMySQL installs MySQL as a system service
func InstallMySQL() error {
	fmt.Println("Installing MySQL service...")

	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base directory: %w", err)
	}

	mysqlDir := filepath.Join(baseDir, "mysql")
	dataDir := filepath.Join(mysqlDir, "data")

	// Check if MySQL is initialized
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return fmt.Errorf("MySQL not initialized. Run 'database mysql init' first")
	}

	// Generate MySQL config
	if err := generateMySQLConfig(); err != nil {
		return fmt.Errorf("failed to generate MySQL config: %w", err)
	}

	// Create service
	svc, err := createMySQLService()
	if err != nil {
		return err
	}

	// Install service
	if err := svc.Install(); err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}

	fmt.Println("✅ MySQL service installed successfully")
	fmt.Printf("   Service Name: %s\n", mysqlServiceName)
	fmt.Println("   Use 'database mysql start' to start MySQL")
	return nil
}

// UninstallMySQL uninstalls MySQL system service
func UninstallMySQL() error {
	fmt.Println("Uninstalling MySQL service...")

	svc, err := createMySQLService()
	if err != nil {
		// If service doesn't exist, consider it already uninstalled
		fmt.Println("⚠️  MySQL service not found (may already be uninstalled)")
		return nil
	}

	// Check status
	status, err := svc.Status()
	if err == nil && status == service.StatusRunning {
		fmt.Println("Stopping MySQL service first...")
		if err := svc.Stop(); err != nil {
			fmt.Printf("Warning: failed to stop service: %v\n", err)
		}
	}

	// Uninstall service
	if err := svc.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}

	fmt.Println("✅ MySQL service uninstalled successfully")
	return nil
}

// StartMySQL starts the MySQL service
func StartMySQL() error {
	fmt.Println("Starting MySQL service...")

	svc, err := createMySQLService()
	if err != nil {
		return err
	}

	// Check if already running
	status, err := svc.Status()
	if err == nil && status == service.StatusRunning {
		fmt.Println("MySQL service is already running")
		return nil
	}

	// Start service
	if err := svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	fmt.Println("✅ MySQL service started successfully")
	return nil
}

// StopMySQL stops the MySQL service
func StopMySQL() error {
	fmt.Println("Stopping MySQL service...")

	svc, err := createMySQLService()
	if err != nil {
		return err
	}

	// Check if running
	status, err := svc.Status()
	if err == nil && status == service.StatusStopped {
		fmt.Println("MySQL service is not running")
		return nil
	}

	// Stop service
	if err := svc.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	fmt.Println("✅ MySQL service stopped successfully")
	return nil
}

// StatusMySQL checks MySQL service status
func StatusMySQL() error {
	svc, err := createMySQLService()
	if err != nil {
		return err
	}

	status, err := svc.Status()
	if err != nil {
		return fmt.Errorf("failed to get service status: %w", err)
	}

	switch status {
	case service.StatusRunning:
		fmt.Println("MySQL service is running")
	case service.StatusStopped:
		fmt.Println("MySQL service is stopped")
	default:
		fmt.Printf("MySQL service status: %v\n", status)
	}

	return nil
}

// RestartMySQL restarts the MySQL service
func RestartMySQL() error {
	fmt.Println("Restarting MySQL service...")

	svc, err := createMySQLService()
	if err != nil {
		return err
	}

	// Check status
	status, err := svc.Status()
	if err == nil && status == service.StatusRunning {
		fmt.Println("Stopping MySQL service...")
		if err := svc.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
	}

	// Start service
	fmt.Println("Starting MySQL service...")
	if err := svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	fmt.Println("✅ MySQL service restarted successfully")
	return nil
}
