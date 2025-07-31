# GoForge AI Agent Instructions

GoForge is a Go CLI template/framework for creating production-ready command-line applications with automatic versioning, multi-platform builds, and comprehensive logging.

## Architecture Overview

**Dual Interface Pattern**: GoForge functions both as a standalone CLI and as a library:
- CLI mode: Entry via `cmd/main.go` → `RegX().Command().Execute()`
- Library mode: Implement the `GoForge` interface in `goforge.go`

**Core Components**:
- `cmd/wrpr.go`: Main command wrapper implementing the GoForge interface
- `cmd/cli/`: Command definitions and CLI utilities
- `logger/`: Structured logging with context and colors using external `logz` package
- `version/`: Automatic semantic versioning with GitHub integration

## Key Patterns

**Command Registration**: Add new commands in `cmd/cli/service.go` via `ServiceCmdList()`, then register in `wrpr.go`:
```go
rtCmd.AddCommand(cc.ServiceCmdList()...)
```

**Logger Usage**: Import as `gl "github.com/rafa-mori/goforge/logger"` and use structured logging:
```go
gl.Log("info", "message")
gl.Log("debug", map[string]interface{}{"key": "value"})
```

**Banner System**: Random ASCII banners in `cmd/cli/common.go`, controlled by `GOFORGE_PRINT_BANNER` env var.

## Build System

**Primary Commands**:
- `make build`: Builds binary via `support/install.sh build`
- `make install`: Full installation via `support/install.sh install`

**CI/CD Flow**: GitHub Actions in `.github/workflows/release.yml`:
1. Builds for Linux/amd64 with ldflags injection
2. UPX compression for size optimization
3. Auto-publishes to GitHub Releases with checksums
4. Version injected into `version/CLI_VERSION` file

**Multi-platform**: Build script supports Linux, macOS, Windows with different architectures via `support/build.sh`.

## Developer Conventions

**File Organization**: Follow the established pattern:
- `cmd/main.go`: CLI entry point
- `cmd/cli/*.go`: Individual command implementations
- `goforge.go`: Library interface definition
- `support/`: Build and installation scripts

**Coding Standards**: Follow `support/instructions/go.md` - table-driven tests, godoc comments, composition over inheritance, explicit error handling.

**Environment Variables**:
- `GOFORGE_PRINT_BANNER`: Controls banner display
- `GITHUB_OWNER`: Overrides project owner for version checks
- `MODULE_ALIAS`: Sets module alias for logging

**Versioning**: The `version/` package handles automatic version detection from GitHub releases and git tags, with fallback to embedded version from CI/CD.

## Integration Points

- External logging via `github.com/rafa-mori/logz`
- CLI framework via `github.com/spf13/cobra`
- GitHub API for version checking and releases
- UPX for binary compression in builds

When adding features, maintain the dual CLI/library interface and ensure commands follow the established registration pattern.
