# Distribution & Versioning Flow

![GoForge Distribution](assets/top_banner_m_a.png)

---

[![Build](https://github.com/rafa-mori/goforge/actions/workflows/release.yml/badge.svg)](https://github.com/rafa-mori/goforge/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.20-blue)](../go.mod)
[![Releases](https://img.shields.io/github/v/release/rafa-mori/goforge?include_prereleases)](https://github.com/rafa-mori/goforge/releases)

---

[🇧🇷 Leia esta documentação em Português](DISTRIBUTION.pt-BR.md) | [← Back to README](../README.md)

This document describes GoForge's advanced distribution and versioning system, which provides centralized application information management and automated binary generation.

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Manifest-Driven Architecture](#-manifest-driven-architecture)
- [Application Info Package](#-application-info-package)
- [Binary Generation](#️-binary-generation)
- [Versioning System](#-versioning-system)
- [Build Artifacts](#-build-artifacts)
- [CI/CD Integration](#-cicd-integration)
- [Usage Examples](#-usage-examples)

---

## 🎯 Overview

GoForge features a sophisticated distribution system that centralizes all application metadata in a single source of truth: the `info/manifest.json` file. This approach ensures consistency across builds, simplifies configuration management, and enables powerful automation capabilities.

### Key Features

- **Centralized Configuration**: Single `manifest.json` file controls all application metadata
- **Embedded Manifest**: Application info is embedded at compile time for runtime access
- **Automated Versioning**: Semantic versioning with GitHub integration
- **Multi-platform Builds**: Support for Linux, macOS, and Windows
- **Binary Optimization**: Automatic UPX compression for smaller binaries
- **Organized Artifacts**: All build outputs go to the dedicated `bin/` directory

---

## 📄 Manifest-Driven Architecture

### The Central Manifest

The `info/manifest.json` file serves as the single source of truth for all application metadata:

```json
{
  "name": "GoForge",
  "application": "goforge",
  "version": "1.0.1",
  "private": false,
  "published": true,
  "aliases": ["goforge"],
  "repository": "https://github.com/rafa-mori/goforge",
  "homepage": "https://github.com/rafa-mori/goforge",
  "description": "A Go scaffolding tool for creating projects with a standard structure and best practices.",
  "main": "cmd/main.go",
  "bin": "goforge",
  "author": "Rafael Mori <faelmori@gmail.com>",
  "organization": "rafa-mori",
  "license": "MIT",
  "keywords": ["goforge", "scaffolding", "tool", "management"],
  "platforms": ["linux"]
}
```

### Manifest Fields

| Field | Description | Example |
|-------|-------------|---------|
| `name` | Display name of the application | `"GoForge"` |
| `application` | Internal application name | `"goforge"` |
| `version` | Current version (semver) | `"1.0.1"` |
| `bin` | Binary name when built | `"goforge"` |
| `repository` | Git repository URL | `"https://github.com/..."` |
| `description` | Short description | `"A Go scaffolding tool..."` |
| `main` | Entry point file | `"cmd/main.go"` |
| `author` | Author information | `"Name <email>"` |
| `organization` | Organization/owner | `"rafa-mori"` |
| `platforms` | Supported platforms | `["linux", "darwin", "windows"]` |
| `private` | Private repository flag | `false` |
| `published` | Published status | `true` |

---

## 🔧 Application Info Package

### Package Structure

The `info` package provides type-safe access to the manifest data:

```go
// info/application.go
package manifest

type Manifest interface {
    GetName() string
    GetVersion() string
    GetRepository() string
    GetDescription() string
    // ... other getters
}
```

### Embedded Manifest

The manifest is embedded at compile time using Go's `embed` directive:

```go
//go:embed manifest.json
var manifestJSONData []byte
```

### Runtime Access

Applications can access manifest data at runtime:

```go
import manifest "github.com/rafa-mori/goforge/info"

func main() {
    info, err := manifest.GetManifest()
    if err != nil {
        log.Fatal("Failed to get manifest:", err)
    }
    
    fmt.Printf("Application: %s v%s\n", 
        info.GetName(), 
        info.GetVersion())
}
```

---

## 🏗️ Binary Generation

### Build Directory Structure

All build artifacts are organized in the `bin/` directory:

```text
bin/
├── goforge_linux_amd64           # Linux binary
├── goforge_linux_amd64.tar.gz    # Compressed archive
├── goforge_darwin_amd64          # macOS binary
├── goforge_darwin_amd64.tar.gz   # macOS archive
├── goforge_windows_amd64.exe     # Windows binary
└── goforge_windows_amd64.zip     # Windows archive
```

### Build Process

1. **Manifest Reading**: Build scripts read `info/manifest.json` for configuration
2. **Binary Compilation**: Go build with embedded manifest data
3. **Platform Targeting**: Cross-compilation for multiple platforms
4. **Optimization**: UPX compression (when enabled)
5. **Packaging**: Archive creation with checksums

### Makefile Integration

The Makefile automatically reads manifest data:

```makefile
APP_NAME := $(shell jq -r '.name' < $(ROOT_DIR)info/manifest.json)
BINARY_NAME := $(shell jq -r '.bin' < $(ROOT_DIR)info/manifest.json)
ORGANIZATION := $(shell jq -r '.organization' < $(ROOT_DIR)info/manifest.json)
```

---

## 🔄 Versioning System

### Version Management

The versioning system provides:

- **Current Version**: Read from `info/manifest.json`
- **Latest Version**: Fetched from GitHub releases API
- **Version Comparison**: Semantic version parsing and comparison
- **Update Checks**: Automatic checks for newer versions

### Version Service Interface

```go
type Service interface {
    GetLatestVersion() (string, error)
    GetCurrentVersion() string
    IsLatestVersion() (bool, error)
    GetRepository() string
}
```

### CLI Integration

```bash
# Check current version
goforge version

# Check for updates
goforge version --check

# Get version info
goforge version --info
```

---

## 📁 Build Artifacts

### Local Development

```bash
# Build for current platform
make build

# Install locally
make install

# Clean build artifacts
make clean
```

### Multi-platform Build

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build PLATFORM=linux ARCH=amd64
```

### Artifact Types

- **Raw Binaries**: `goforge_platform_arch[.exe]`
- **Compressed Archives**: `.tar.gz` (Unix) / `.zip` (Windows)
- **Checksums**: SHA256 hashes for integrity verification
- **Debug Symbols**: Separate debug information (when enabled)

---

## 🚀 CI/CD Integration

### GitHub Actions Workflow

The automated workflow:

1. **Checkout Code**: Get source code and manifest
2. **Setup Environment**: Configure Go and build tools
3. **Read Manifest**: Extract build configuration
4. **Multi-platform Build**: Generate binaries for all platforms
5. **Optimization**: Apply UPX compression
6. **Archive Creation**: Package binaries with documentation
7. **Checksum Generation**: Create integrity hashes
8. **Release Publishing**: Upload to GitHub releases

### Environment Variables

| Variable | Description | Source |
|----------|-------------|--------|
| `APP_NAME` | Application name | `manifest.json` |
| `VERSION` | Current version | `manifest.json` |
| `REPOSITORY` | Repository URL | `manifest.json` |
| `BINARY_NAME` | Output binary name | `manifest.json` |

---

## 💡 Usage Examples

### 1. Creating a New Release

1. Update version in `info/manifest.json`:

   ```json
   {
     "version": "1.0.2"
   }
   ```

1. Commit and tag:

   ```bash
   git add info/manifest.json
   git commit -m "bump: version 1.0.2"
   git tag v1.0.2
   git push origin main --tags
   ```

1. CI/CD automatically builds and releases

### 2. Adding Platform Support

Update the platforms array in `manifest.json`:

```json
{
  "platforms": ["linux", "darwin", "windows"]
}
```

### 3. Customizing Build Output

Modify the binary name:

```json
{
  "bin": "my-custom-name"
}
```

This generates `my-custom-name_linux_amd64` instead of `goforge_linux_amd64`.

### 4. Runtime Manifest Access

```go
package main

import (
    "fmt"
    manifest "github.com/rafa-mori/goforge/info"
)

func main() {
    info, _ := manifest.GetManifest()
    
    fmt.Printf("Welcome to %s v%s\n", 
        info.GetName(), 
        info.GetVersion())
    fmt.Printf("Repository: %s\n", 
        info.GetRepository())
    fmt.Printf("Author: %s\n", 
        info.GetAuthor())
}
```

---

## 🎯 Benefits

### For Developers

- **Single Source of Truth**: No more scattered configuration files
- **Type Safety**: Compile-time validation of manifest access
- **Automated Workflows**: Hands-off build and release process
- **Consistent Binaries**: Same configuration across all builds

### For Users

- **Organized Downloads**: Clear binary naming and organization
- **Integrity Verification**: Checksums for security
- **Multi-platform Support**: Native binaries for all platforms
- **Optimized Size**: UPX compression reduces download time

### For Operations

- **Reproducible Builds**: Deterministic build process
- **Audit Trail**: Full traceability of build configurations
- **Easy Deployment**: Standardized artifact structure
- **Version Management**: Automated version tracking and comparison

---

## 🤝 Contributing

When contributing to GoForge's distribution system:

1. **Test Locally**: Run `make build` to verify changes
2. **Update Manifest**: Modify `info/manifest.json` as needed
3. **Document Changes**: Update this documentation
4. **Test CI/CD**: Verify GitHub Actions still work

---

## 📄 License

MIT License - see [LICENSE](../LICENSE) file for details.

---

## 👤 Author

**Rafael Mori** - [@rafa-mori](https://github.com/rafa-mori)

---

> Made with 💙 for streamlined Go development workflows
