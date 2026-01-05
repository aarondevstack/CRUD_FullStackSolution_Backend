# MySQL Assets Directory

This directory should contain the MySQL ZIP packages for each platform.

## Required Files

Download the following files and place them in the respective subdirectories:

### Darwin (macOS)
- **File**: `darwin/mysql-8.4.7-macos15-arm64.zip`
- **Download**: [mysql-8.4.7-macos15-arm64.zip](https://github.com/aarondevstack/mysql-assets/blob/main/darwin/mysql-8.4.7-macos15-arm64.zip)
- **Size**: ~400MB

### Linux
- **File**: `linux/mysql-8.4.7-linux-x86_64.zip`
- **Download**: [mysql-8.4.7-linux-x86_64.zip](https://github.com/aarondevstack/mysql-assets/blob/main/linux/mysql-8.4.7-linux-x86_64.zip)
- **Size**: ~400MB

### Windows
- **File**: `windows/mysql-8.4.7-winx64.zip`
- **Download**: [mysql-8.4.7-winx64.zip](https://github.com/aarondevstack/mysql-assets/blob/main/windows/mysql-8.4.7-winx64.zip)
- **Size**: ~400MB

## Directory Structure

```
assets/
├── darwin/
│   └── mysql-8.4.7-macos15-arm64.zip
├── linux/
│   └── mysql-8.4.7-linux-x86_64.zip
└── windows/
    └── mysql-8.4.7-winx64.zip
```

## Note

These files are embedded into the executable using `go:embed` with platform-specific build tags. Only the ZIP file for the target platform will be included in the final binary.

**Important**: These files are excluded from git via `.gitignore` due to their large size.
