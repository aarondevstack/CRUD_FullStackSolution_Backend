# CRUD Solution Backend

This is the backend service for the CRUD Fullstack Solution, built with **Golang**. It follows a modular architecture using **Cobra** for CLI, **GoFiber** for the Web API, and **Ent** for the Data Layer.

## Features

- **Modular CLI**: Built with Cobra, supporting clear separation of command definition and implementation.
- **RESTful API**: Powered by GoFiber, featuring JWT authentication and Casbin RBAC.
- **ORM & Migrations**: Uses Ent for strict schema definition and Atlas for database migrations.
- **Cross-Platform Service Management**: Built-in commands to install, start, stop, and manage services on macOS, Linux, and Windows.
- **Embedded Database**: Supports embedding MySQL binaries for a self-contained deployment.

## Prerequisites

- **Go**: 1.23 or higher
- **Make**: For build scripts
- **Atlas**: For local schema generation (optional, as binary is embedded for prod)

## Getting Started

### 1. Build the Project

Use the provided `Makefile` to compile for your target platform:

```bash
# macOS
make build-darwin

# Linux
make build-linux

# Windows
make build-windows
```

### 2. Run in Development (Hot Reload)

Ensure you have `air` installed for hot reloading:

```bash
air
```

This will start the API service in development mode, listening on port `8888`.

### 3. Service Management

The application binary includes built-in commands to manage itself and the database as system services.

**MySQL Service:**
```bash
# Install as a service
sudo ./bin/crud-solution services mysql install

# Start the service
sudo ./bin/crud-solution services mysql start

# Check status
sudo ./bin/crud-solution services mysql status
```

**API Service:**
```bash
# Install as a service
sudo ./bin/crud-solution services api install

# Start the service
sudo ./bin/crud-solution services api start
```

## Directory Structure

- `cmd/`: Command definitions (Cobra).
- `internal/cli/`: Command implementations.
- `internal/services/api`: GoFiber Web Service logic.
- `internal/database/ent`: Ent ORM schema and generated code.
- `internal/database/migrations`: Atlas migration files.

