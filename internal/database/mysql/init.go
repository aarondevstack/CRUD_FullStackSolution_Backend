package mysql

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/config"
)

// InitMySQL initializes MySQL database
func InitMySQL() error {
	fmt.Println("Initializing MySQL...")

	// Get base directory
	baseDir, err := config.GetBaseDir()
	if err != nil {
		return fmt.Errorf("failed to get base directory: %w", err)
	}

	mysqlDir := filepath.Join(baseDir, "mysql")
	dataDir := filepath.Join(mysqlDir, "data")

	// Check if already initialized
	if _, err := os.Stat(dataDir); err == nil {
		fmt.Println("MySQL already initialized (data directory exists)")
		return nil
	}

	// Step 1: Decompress MySQL ZIP
	fmt.Println("Step 1/4: Decompressing MySQL ZIP package...")
	if err := decompressMySQLZip(mysqlDir); err != nil {
		return fmt.Errorf("failed to decompress MySQL: %w", err)
	}

	// Step 2: Initialize data directory
	fmt.Println("Step 2/4: Initializing MySQL data directory...")
	mysqldPath := filepath.Join(mysqlDir, "bin", "mysqld")
	cmd := exec.Command(mysqldPath, "--initialize-insecure", fmt.Sprintf("--datadir=%s", dataDir))
	// Set library path for macOS
	cmd.Env = append(os.Environ(), fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(mysqlDir, "lib")))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize data directory: %w\nOutput: %s", err, string(output))
	}

	// Step 3: Generate and execute init SQL
	fmt.Println("Step 3/4: Setting up database and users...")
	if err := executeInitSQL(mysqlDir, dataDir); err != nil {
		return fmt.Errorf("failed to execute init SQL: %w", err)
	}

	// Step 4: Clean up
	fmt.Println("Step 4/4: Cleaning up temporary files...")

	fmt.Println("✅ MySQL initialization completed successfully!")
	return nil
}

// decompressMySQLZip decompresses the embedded MySQL ZIP to the target directory
// Handles two scenarios:
// 1. With top-level directory: Rename it to "mysql"
// 2. Without top-level directory: Create "mysql" and move all contents
func decompressMySQLZip(targetDir string) error {
	// Create temporary extraction directory
	baseDir := filepath.Dir(targetDir)
	tempDir := filepath.Join(baseDir, "mysql_temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Read ZIP from embedded bytes
	zipReader, err := zip.NewReader(bytes.NewReader(mysqlZip), int64(len(mysqlZip)))
	if err != nil {
		return fmt.Errorf("failed to read ZIP: %w", err)
	}

	// Extract all files to temp directory
	for _, file := range zipReader.File {
		targetPath := filepath.Join(tempDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, file.Mode())
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file in ZIP: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	// Check if there's a single top-level directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	if len(entries) == 1 && entries[0].IsDir() {
		// Scenario 1: Single top-level directory - rename it to "mysql"
		topLevelDir := filepath.Join(tempDir, entries[0].Name())
		fmt.Printf("Detected top-level directory: %s (renaming to mysql)\n", entries[0].Name())
		if err := os.Rename(topLevelDir, targetDir); err != nil {
			return fmt.Errorf("failed to rename directory: %w", err)
		}
	} else {
		// Scenario 2: Flat structure - create mysql directory and move all contents
		fmt.Println("No single top-level directory detected (flat structure)")
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		// Move all contents to target directory
		for _, entry := range entries {
			srcPath := filepath.Join(tempDir, entry.Name())
			dstPath := filepath.Join(targetDir, entry.Name())
			if err := os.Rename(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to move %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// executeInitSQL generates and executes initialization SQL
func executeInitSQL(mysqlDir, dataDir string) error {
	// Generate init SQL
	initSQL := `
-- Set root password
ALTER USER 'root'@'localhost' IDENTIFIED BY 'aarondevstack@2026';
FLUSH PRIVILEGES;

-- Create database
CREATE DATABASE IF NOT EXISTS ` + "`crud-db`" + ` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create crud user
CREATE USER IF NOT EXISTS 'crud'@'localhost' IDENTIFIED BY 'crud@2026';
CREATE USER IF NOT EXISTS 'crud'@'%' IDENTIFIED BY 'crud@2026';

-- Grant privileges
GRANT ALL PRIVILEGES ON ` + "`crud-db`" + `.* TO 'crud'@'localhost';
GRANT ALL PRIVILEGES ON ` + "`crud-db`" + `.* TO 'crud'@'%';
FLUSH PRIVILEGES;

-- Shutdown server gracefully
SHUTDOWN;
`

	// Write init SQL to temporary file
	initSQLPath := filepath.Join(mysqlDir, "mysql_init.sql")
	if err := os.WriteFile(initSQLPath, []byte(initSQL), 0644); err != nil {
		return fmt.Errorf("failed to write init SQL: %w", err)
	}
	defer os.Remove(initSQLPath)

	// Start MySQL temporarily with init-file
	mysqldPath := filepath.Join(mysqlDir, "bin", "mysqld")
	cmd := exec.Command(mysqldPath,
		fmt.Sprintf("--init-file=%s", initSQLPath),
		fmt.Sprintf("--datadir=%s", dataDir),
	)
	// Set library path for macOS
	cmd.Env = append(os.Environ(), fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(mysqlDir, "lib")))

	// Run the command and wait for it to complete
	// MySQL will shutdown gracefully after executing the init SQL (which includes SHUTDOWN command)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to execute init SQL: %w\nOutput: %s", err, string(output))
	}

	return nil
}
