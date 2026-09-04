# gdrive

Google Drive CLI client - A read-only command-line tool for listing folder contents and inspecting file metadata in Google Drive.

## Installation

### Using Go

```bash
go install github.com/longkey1/gdrive@latest
```

### Using Homebrew (coming soon)

```bash
brew install longkey1/tap/gdrive
```

## Setup

gdrive supports two authentication methods: OAuth (for personal use) and Service Account (for automation/CI).

### OAuth

#### 1. Create Google Cloud Project and OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Drive API
4. Create OAuth 2.0 credentials (Desktop application type)
5. Download the credentials JSON file

If you already use gdoc, gsheet, or gcal, the same OAuth client can be reused.

#### 2. Create Configuration File

Create `~/.config/gdrive/config.toml` (a template is provided as `config.toml.dist`):

```toml
auth_type = "oauth"
application_credentials = "/path/to/credentials.json"
user_credentials = "/path/to/token.json"
```

#### 3. Authenticate

```bash
gdrive auth
```

This will open your browser for Google OAuth authentication. gdrive requests only the `drive.metadata.readonly` scope.

### Service Account

#### 1. Create a Service Account

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Drive API
4. Create a Service Account and download its credentials JSON file
5. Share the target folders with the service account's email address

#### 2. Create Configuration File

Create `~/.config/gdrive/config.toml`:

```toml
auth_type = "service_account"
application_credentials = "/path/to/service-account.json"
```

`user_credentials` and the `gdrive auth` command are not needed for Service Account authentication.

## Usage

### List Folder Contents

`files list` returns the direct children of a folder. It is not recursive, excludes trashed items, and works for folders inside shared drives.

```bash
# List by folder ID
gdrive files list --parent <folder-id>

# Pass a Google Drive folder URL instead of an ID
gdrive files list --parent "https://drive.google.com/drive/u/0/folders/<folder-id>"

# Filter by name substring
gdrive files list --parent <folder-id> -q "recording"

# Limit the number of results (default: 100, 0 = unlimited)
gdrive files list --parent <folder-id> --max-results 500

# Output as JSON
gdrive files list --parent <folder-id> --format json
```

JSON output is always an array (`[]` when the folder is empty). Each element has the following fields, all of which are always present:

```json
[
  {
    "id": "1AbC...",
    "name": "2026-09-01 catch-up.mp4",
    "mimeType": "video/mp4",
    "webViewLink": "https://drive.google.com/file/d/1AbC.../view?usp=drivesdk",
    "modifiedTime": "2026-09-01T10:00:00.000Z",
    "size": 123456789
  }
]
```

Example: produce `mimeType<TAB>webViewLink<TAB>name` lines for a script.

```bash
gdrive files list --parent <folder-id> --format json \
  | jq -r '(. // [])[] | [.mimeType, .webViewLink, .name] | @tsv'
```

### Get File Metadata

`files get` accepts a file ID or a Google Drive / Docs / Sheets / Slides URL.

```bash
gdrive files get <file-id>
gdrive files get "https://drive.google.com/file/d/<file-id>/view"
gdrive files get "https://docs.google.com/document/d/<doc-id>/edit"
gdrive files get <file-id> --format json
```

JSON output additionally includes `parents` and `owners` (`emailAddress`, `displayName`).

### Accepted URL Formats

| Kind | Example |
|------|---------|
| Folder | `https://drive.google.com/drive/folders/<id>` |
| Folder (account index) | `https://drive.google.com/drive/u/0/folders/<id>` |
| File | `https://drive.google.com/file/d/<id>/view` |
| Docs / Slides / Sheets | `https://docs.google.com/(document\|presentation\|spreadsheets)/d/<id>/edit` |

Trailing query strings (`?usp=sharing`) and fragments are ignored.

### Output Formats

| Format | Description |
|--------|-------------|
| `text` | Table output (default). `files list` shows `ID / NAME / TYPE / MODIFIED` |
| `json` | JSON output (`files list`: array, `files get`: object) |

### Shell Completion

```bash
gdrive completion zsh > "${fpath[1]}/_gdrive"
gdrive completion bash > /etc/bash_completion.d/gdrive
```

### Version

```bash
gdrive version

# Show only the version number
gdrive version --short
```

## Configuration Options

| Option | Description |
|--------|-------------|
| `auth_type` | Authentication type: `oauth` or `service_account` |
| `application_credentials` | Path to OAuth client credentials JSON file |
| `user_credentials` | Path to store OAuth user token (for OAuth auth type) |
| `read_only` | Accepted for parity with gdoc/gsheet/gcal. gdrive has no write commands, so it currently has no effect. Also settable via `GDRIVE_READ_ONLY` and `--read-only` |

## License

MIT
