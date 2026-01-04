package database

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/atlas"
)

// MigrateDatabase applies database migrations using embedded Atlas
func MigrateDatabase() error {
	fmt.Println("Applying database migrations...")

	// Load configuration
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build database URL
	dbURL := fmt.Sprintf("mysql://%s:%s@%s:%d/%s?parseTime=true",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name,
	)

	// Create temporary directory
	tempDir := "crud_temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		// Clean up temp directory
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Printf("Warning: failed to clean up temp directory: %v\n", err)
		}
	}()

	// Extract Atlas binary
	atlasPath, err := atlas.ExtractAtlas(tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract Atlas: %w", err)
	}

	// Extract migrations
	if err := extractMigrations(tempDir); err != nil {
		return fmt.Errorf("failed to extract migrations: %w", err)
	}

	migrationsDir := filepath.Join(tempDir, "migrations")

	// Run Atlas migrate
	if err := atlas.RunMigrate(atlasPath, migrationsDir, dbURL); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✅ Database migrations applied successfully!")
	return nil
}

// extractMigrations extracts all migration files to a temporary directory
func extractMigrations(tempDir string) error {
	migrationsDir := filepath.Join(tempDir, "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Walk through embedded migrations
	return fs.WalkDir(MigrationsFS, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "migrations" directory itself
		if path == "migrations" {
			return nil
		}

		// Remove "migrations/" prefix from path
		relativePath := path[11:] // len("migrations/") = 11

		targetPath := filepath.Join(migrationsDir, relativePath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file content
		content, err := fs.ReadFile(MigrationsFS, path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		// Write to temp directory
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write migration file %s: %w", targetPath, err)
		}

		return nil
	})
}

// BackupDatabase creates a database backup
func BackupDatabase() error {
	fmt.Println("Backing up database...")

	// Load configuration
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base dir: %w", err)
	}

	// Paths
	mysqldumpName := "mysqldump"
	if os.Getenv("GOOS") == "windows" {
		mysqldumpName = "mysqldump.exe"
	}
	mysqldumpPath := filepath.Join(baseDir, "mysql", "bin", mysqldumpName)
	backupDir := filepath.Join(baseDir, "crud_backup")

	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)
	backupPath := filepath.Join(backupDir, filename)

	// Build command
	// mysqldump -u[user] -p[password] -h[host] -P[port] --databases [dbname]
	args := []string{
		fmt.Sprintf("-u%s", config.AppConfig.Database.User),
		fmt.Sprintf("-p%s", config.AppConfig.Database.Password),
		fmt.Sprintf("-h%s", config.AppConfig.Database.Host),
		fmt.Sprintf("-P%d", config.AppConfig.Database.Port),
		"--databases",
		config.AppConfig.Database.Name,
	}

	// Create output file
	outfile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outfile.Close()

	cmd := exec.Command(mysqldumpPath, args...)
	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysqldump failed: %w", err)
	}

	fmt.Printf("✅ Database backup created: %s\n", backupPath)
	return nil
}

// RestoreDatabase restores database from backup
func RestoreDatabase(backupFile string) error {
	fmt.Printf("Restoring database from %s...\n", backupFile)

	// Validate file exists
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFile)
	}

	// Load configuration
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base dir: %w", err)
	}

	// Paths
	mysqlName := "mysql"
	if os.Getenv("GOOS") == "windows" {
		mysqlName = "mysql.exe"
	}
	mysqlPath := filepath.Join(baseDir, "mysql", "bin", mysqlName)

	// Build command
	// mysql -u[user] -p[password] -h[host] -P[port] [dbname] < [backupFile]
	args := []string{
		fmt.Sprintf("-u%s", config.AppConfig.Database.User),
		fmt.Sprintf("-p%s", config.AppConfig.Database.Password),
		fmt.Sprintf("-h%s", config.AppConfig.Database.Host),
		fmt.Sprintf("-P%d", config.AppConfig.Database.Port),
		config.AppConfig.Database.Name,
	}

	// Open input file
	infile, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer infile.Close()

	cmd := exec.Command(mysqlPath, args...)
	cmd.Stdin = infile
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %w", err)
	}

	fmt.Println("✅ Database restored successfully!")
	return nil
}

// DeployDatabase runs the full database deployment flow
func DeployDatabase() error {
	fmt.Println("Database deployment - to be implemented")
	// TODO: Implement: init -> install -> start -> migrate -> seed
	return nil
}

// RestartDatabase restarts the database service
func RestartDatabase() error {
	fmt.Println("Database restart - to be implemented")
	// TODO: Implement database restart
	return nil
}

// StopDatabase stops the database service
func StopDatabase() error {
	fmt.Println("Database stop - to be implemented")
	// TODO: Implement database stop
	return nil
}

// UninstallDatabase uninstalls the database service
func UninstallDatabase() error {
	fmt.Println("Database uninstall - to be implemented")
	// TODO: Implement database uninstall
	return nil
}
