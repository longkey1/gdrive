# AGENTS.md

This file provides guidance to AI coding agents (Claude Code, etc.) when working with code in this repository.

## Project Overview

gdrive is a read-only CLI client for Google Drive. It uses the Drive API v3 to list the contents of a folder and to show file metadata, in plain text or JSON. Built with Go 1.27, cobra, and viper. It is a sister tool of gdoc / gsheet / gcal and shares their layout and conventions.

## Build & Development Commands

```bash
make build       # Build binary to bin/gdrive
make test        # Run tests (go test ./...)
make fmt         # Format code (go fmt ./...)
make vet         # Vet code (go vet ./...)
make lint        # Run golangci-lint (version managed by go.mod tool directive)
make tidy        # Tidy dependencies (go mod tidy)
```

Release: `make release type=patch|minor|major dryrun=false` (GitHub Actions + goreleaser)

## Architecture

```
main.go                     # Entry point. Version info injected via ldflags
cmd/                        # cobra command definitions
  root.go                   # Config loading ($HOME/.config/gdrive/config.toml via viper), --read-only guard
  auth.go                   # OAuth flow (only for auth_type=oauth)
  files.go                  # `files list --parent <id|url>` and `files get <id|url>`
  version.go                # `version` command
internal/
  gdrive/
    config.go               # Config struct (auth_type, application_credentials, user_credentials, read_only)
    service.go              # Service: creates Authenticator based on auth_type, initializes Drive service
    files.go                # ListFiles (paged, shared drives included), GetFile, BuildListQuery
    url.go                  # ParseFileInput: Drive/Docs/Sheets/Slides URL or bare ID → file ID
    format.go               # Output formatters (text via tablewriter, json via MarshalIndent)
  google/
    auth.go                 # Authenticator interface, OAuthAuthenticator, ServiceAccountAuthenticator
    drive.go                # DriveService wrapper
  version/                  # Version info set at build time via ldflags
```

### Key Design Patterns

- **Auth abstraction**: `google.Authenticator` interface switches between OAuth and Service Account. `gdrive.Service` selects the concrete type based on `auth_type` in config.
- **Read-only by design**: v1 has no write commands. The OAuth scope is `drive.metadata.readonly` only. `--read-only` / `read_only` are accepted for parity with the sister tools and guard commands annotated with `write`, of which there are currently none.
- **Stable JSON contract**: `files list --format json` always emits an array (`[]` when empty) whose elements carry `id`, `name`, `mimeType`, `webViewLink`, `modifiedTime`, `size`. Downstream scripts pipe this into `jq`; do not rename or drop these fields.
- **Shared drives**: `files.list` / `files.get` are called with `supportsAllDrives` and `includeItemsFromAllDrives` so folders inside shared drives work.

## Configuration

Config file: `$HOME/.config/gdrive/config.toml` (see `config.toml.dist`)

```toml
auth_type = "oauth"                    # "oauth" or "service_account"
application_credentials = "/path/to/credentials.json"
user_credentials = "/path/to/token.json"
```

Environment variables use the `GDRIVE_` prefix (e.g. `GDRIVE_READ_ONLY=true`).

## Environment

- direnv (`.envrc`) sets GOROOT/GOPATH/GOBIN/GOCACHE
- Go version managed via `got` command (`got path 1.27`)
- Product name stored in `.product_name` file
