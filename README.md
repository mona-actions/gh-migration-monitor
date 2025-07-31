# GitHub Migration Monitor

A [GitHub CLI](https://cli.github.com) extension that monitors the progress of GitHub Organization migrations with a real-time terminal UI dashboard.

This extension provides a comprehensive view of repository migrations, supporting both legacy migrations and the new GitHub Enterprise Importer (GEI) migrations.

## Features

- **Real-time Monitoring**: Live dashboard that updates migration status every 60 seconds
- **Multi-State Tracking**: Monitors migrations across Queued, In Progress, Succeeded, and Failed states
- **Legacy Support**: Compatible with both legacy migrations and GEI migrations
- **Interactive UI**: Keyboard navigation and focus switching between migration states
- **Error Details**: Detailed failure reasons for failed migrations

## Installation

Install the extension using the GitHub CLI:

```bash
gh extension install mona-actions/gh-migration-monitor
```

## Usage

### Basic Usage

Monitor migrations for a GitHub organization:

```bash
gh migration-monitor --organization myorg
```

### Advanced Usage

```bash
# Monitor with a specific GitHub token
gh migration-monitor --organization myorg --github-token ghp_xxxxxxxxxxxx

# Monitor legacy migrations
gh migration-monitor --organization myorg --legacy

# Short flags
gh migration-monitor -o myorg -t ghp_xxxxxxxxxxxx -l
```

### Command Line Options

| Flag             | Short | Description                    | Required |
| ---------------- | ----- | ------------------------------ | -------- |
| `--organization` | `-o`  | GitHub organization to monitor | ✅       |
| `--github-token` | `-t`  | GitHub personal access token   | ❌\*     |
| `--legacy`       | `-l`  | Monitor legacy migrations      | ❌       |

\*The GitHub token can also be provided via the `GHMM_GITHUB_TOKEN` environment variable.

## Configuration

The extension supports configuration through multiple sources:

### Environment Variables

Set these environment variables to avoid passing flags:

```bash
export GHMM_GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GHMM_GITHUB_ORGANIZATION="myorg"
export GHMM_ISLEGACY="true"  # for legacy migrations
```

### Configuration File

Create a configuration file at `~/.gh-migration-monitor/config.yaml`:

```yaml
github:
  token: 'ghp_xxxxxxxxxxxx'
  organization: 'myorg'

migration:
  is_legacy: false

output:
  format: 'table'
  quiet: false
```

## Interactive Controls

Once the dashboard is running, use these keyboard shortcuts:

| Key | Action                              |
| --- | ----------------------------------- |
| `q` | Focus on **Queued** migrations      |
| `i` | Focus on **In Progress** migrations |
| `s` | Focus on **Succeeded** migrations   |
| `f` | Focus on **Failed** migrations      |
| `x` | Exit the application                |

## Migration States

The dashboard categorizes migrations into four states:

### 🟡 Queued

- `QUEUED`: Migration is waiting to start
- `WAITING`: Migration is in the queue

### 🔵 In Progress

- `IN_PROGRESS`: Migration is actively running
- `PREPARING`: Migration is being prepared
- `PENDING`: Migration is pending
- `MAPPING`: Migration is mapping resources
- `ARCHIVE_UPLOADED`: Archive has been uploaded
- `CONFLICTS`: Migration has conflicts to resolve
- `READY`: Migration is ready to proceed
- `IMPORTING`: Migration is importing data

### 🟢 Succeeded

- `SUCCEEDED`: Migration completed successfully
- `UNLOCKED`: Repository has been unlocked
- `IMPORTED`: Migration has been imported

### 🔴 Failed

- `FAILED`: Migration failed
- `FAILED_IMPORT`: Import process failed

## Prerequisites

- [GitHub CLI](https://cli.github.com) installed and authenticated
- GitHub personal access token with appropriate permissions:
  - For GEI migrations: `read:org`, `repo`
  - For legacy migrations: `read:org`, `repo`, `admin:org`

## Troubleshooting

### Common Issues

**"GitHub token is required"**

- Ensure you have set `GHMM_GITHUB_TOKEN` environment variable or passed `--github-token` flag
- Verify your token has the required permissions

**"Organization not found"**

- Check the organization name is correct
- Ensure your token has access to the organization

**"API rate limit exceeded"**

- The extension includes built-in rate limiting, but if you hit limits, wait a few minutes

### Debugging

Enable verbose logging by setting:

```bash
export GHMM_DEBUG=true
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [MIT License](./LICENSE).

---

**Note**: This tool is inspired by [github-migration-monitor](https://github.com/timrogers/github-migration-monitor) and provides enhanced functionality with a modern terminal UI.
