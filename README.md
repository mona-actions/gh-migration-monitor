# GitHub Migration Monitor

A [GitHub CLI](https://cli.github.com) extension that monitors GitHub Organization migrations with a real-time terminal dashboard.

## Quick Start

```bash
# Install
gh extension install mona-actions/gh-migration-monitor

# Run
gh migration-monitor --organization myorg
```

## Features

- 🔄 **Real-time monitoring** with live dashboard updates
- 📊 **Multi-state tracking** (Queued, In Progress, Succeeded, Failed)
- 🔧 **Legacy support** for both GEI and legacy migrations
- ⌨️ **Interactive UI** with keyboard navigation
- 🚨 **Error details** for failed migrations

## Usage

```bash
# Basic usage
gh migration-monitor --organization myorg

# With custom token
gh migration-monitor --organization myorg --github-token ghp_xxxxxxxxxxxx

# Monitor legacy migrations
gh migration-monitor --organization myorg --legacy
```

### Options

| Flag             | Short | Description               | Required |
| ---------------- | ----- | ------------------------- | -------- |
| `--organization` | `-o`  | GitHub organization       | Yes      |
| `--github-token` | `-t`  | GitHub token              | No*      |
| `--legacy`       | `-l`  | Monitor legacy migrations | No       |

*Can use `GHMM_GITHUB_TOKEN` environment variable instead.

## Configuration

### Environment Variables
```bash
export GHMM_GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GHMM_GITHUB_ORGANIZATION="myorg"
export GHMM_ISLEGACY="true"  # for legacy migrations
```

### Config File
Create `~/.gh-migration-monitor/config.yaml`:
```yaml
github:
  token: 'ghp_xxxxxxxxxxxx'
  organization: 'myorg'
migration:
  is_legacy: false
```

## Controls

| Key | Action            |
| --- | ----------------- |
| `q` | Focus Queued      |
| `i` | Focus In Progress |
| `s` | Focus Succeeded   |
| `f` | Focus Failed      |
| `x` | Exit              |

## Migration States

- 🟡 **Queued**: `QUEUED`, `WAITING`
- 🔵 **In Progress**: `IN_PROGRESS`, `PREPARING`, `PENDING`, `MAPPING`, `IMPORTING`, etc.
- 🟢 **Succeeded**: `SUCCEEDED`, `UNLOCKED`, `IMPORTED`
- 🔴 **Failed**: `FAILED`, `FAILED_IMPORT`

## Requirements

- [GitHub CLI](https://cli.github.com) installed and authenticated
- GitHub token with permissions:
  - **GEI migrations**: `read:org`, `repo`
  - **Legacy migrations**: `read:org`, `repo`, `admin:org`

## Troubleshooting

**Token required**: Set `GHMM_GITHUB_TOKEN` or use `--github-token` flag

**Organization not found**: Check organization name and token permissions

**Rate limit**: Wait a few minutes or enable debug mode:
```bash
export GHMM_DEBUG=true
```

## Contributing

Contributions are welcome! Here's how to get started:

### Development Setup
```bash
# Fork and clone
gh repo fork mona-actions/gh-migration-monitor --clone
cd gh-migration-monitor

# Build and test
go mod download
go build -o gh-migration-monitor
gh extension install .
```

### Making Changes
1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes and add tests
3. Submit a pull request

### Project Structure
- `cmd/` - CLI commands (Cobra)
- `internal/api/` - GitHub API clients
- `internal/config/` - Configuration (Viper)
- `internal/services/` - Business logic
- `internal/ui/` - Terminal UI (tview)

## License

This project is licensed under the [MIT License](./LICENSE).

---

**Note**: This tool is inspired by [github-migration-monitor](https://github.com/timrogers/github-migration-monitor) and provides enhanced functionality with a modern terminal UI.
