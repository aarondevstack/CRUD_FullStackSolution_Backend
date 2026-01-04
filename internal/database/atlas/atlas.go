package atlas

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GetAtlasBinary returns the embedded Atlas binary for the current platform
func GetAtlasBinary() []byte {
	return atlasBinary
}

// ExtractAtlas extracts the Atlas binary to a temporary directory and returns the path
func ExtractAtlas(tempDir string) (string, error) {
	// Determine binary name based on OS
	binaryName := "atlas"
	if runtime.GOOS == "windows" {
		binaryName = "atlas.exe"
	}

	atlasPath := filepath.Join(tempDir, binaryName)

	// Write binary to temp directory
	if err := os.WriteFile(atlasPath, atlasBinary, 0755); err != nil {
		return "", fmt.Errorf("failed to write Atlas binary: %w", err)
	}

	return atlasPath, nil
}

// RunMigrate executes Atlas migrate apply command
func RunMigrate(atlasPath, migrationsDir, dbURL string) error {
	cmd := exec.Command(atlasPath, "migrate", "apply",
		"--dir", fmt.Sprintf("file://%s", migrationsDir),
		"--url", dbURL,
	)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("atlas migrate failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}
