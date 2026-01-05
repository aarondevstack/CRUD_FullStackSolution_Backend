# Atlas Database Client Directory

This directory contains the platform-specific **Atlas** binaries used for database migrations. These binaries are embedded into the final application executable.

## Required Files

Download the Atlas binary for each supported platform and place them in the corresponding subdirectories:

### Darwin (macOS)
- **File**: `darwin/atlas-darwin-arm64-latest`
- **Download**: [atlas-darwin-arm64-latest](https://release.ariga.io/atlas/atlas-darwin-arm64-latest)

### Linux
- **File**: `linux/atlas-linux-amd64-latest`
- **Download**: [atlas-linux-amd64-latest](https://release.ariga.io/atlas/atlas-linux-amd64-latest)

### Windows
- **File**: `windows/atlas-windows-amd64-latest.exe`
- **Download**: [atlas-windows-amd64-latest.exe](https://release.ariga.io/atlas/atlas-windows-amd64-latest.exe)

## Directory Structure

Ensure the directory structure is exactly as follows for the embed directives to work:

```
internal/database/atlas/
├── darwin/
│   └── atlas-darwin-arm64-latest
├── linux/
│   └── atlas-linux-amd64-latest
└── windows/
    └── atlas-windows-amd64-latest.exe
```

## Note

These binary files are excluded from the git repository via `.gitignore` to keep the repo size manageable. They are only needed at build time.
