# GitHub Migration Monitor

A [GitHub CLI](https://cli.github.com) extension that monitors GitHub Organization migrations with a real-time terminal dashboard.

## Quick Start

```bash
# Install
gh extension install mona-actions/gh-migration-monitor

# Run
gh migration-monitor --organization myorg
```

### What You'll See
```
┌ Migration Status - myorg ─────────────────────────────────────────┐
│ Repository Name    Migration ID     Status        Created At      │
│ ─────────────────────────────────────────────────────────────────│
│ frontend-app      mig_123456       ✅ SUCCEEDED   2025-11-05 14:30│
│ backend-api       mig_123457       🔄 IN_PROGRESS 2025-11-05 14:25│
│ mobile-app        mig_123458       📋 QUEUED      2025-11-05 14:20│
│ legacy-system     mig_123459       ❌ FAILED      2025-11-05 14:15│
└─────────────────────────────────────────────────────────────────┘
Commands: r Refresh  / Search  x Exit  Filters: a All  q Queued  i In Progress  s Succeeded  f Failed
                                                    Last updated: 14:35:42
```

## Features

- 🔄 **Real-time monitoring** with automatic 30-second refresh intervals
- 📊 **Multi-state tracking** (Queued, In Progress, Succeeded, Failed)
- 🔍 **Advanced filtering** with status-based views and search functionality
- 🎯 **Live search** with real-time repository name filtering
- 📋 **Comprehensive table** showing Repository Name, Migration ID, Status, and Created At
- 🔧 **Legacy support** for both GEI and legacy migrations
- ⌨️ **Interactive UI** with intuitive keyboard navigation
- 🎨 **Color-coded status** indicators for quick visual assessment
- ⚡ **Optimized performance** with efficient data filtering and updates

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
output:
  format: 'table'      # Output format (reserved for future use)
  quiet: false         # Quiet mode (reserved for future use)
```

## Controls

### Navigation & Actions
| Key | Action            |
| --- | ----------------- |
| `r` | Refresh data      |
| `/` | Open search modal |
| `x` | Exit application  |

### Status Filters
| Key | Filter              |
| --- | ------------------- |
| `a` | Show All migrations |
| `q` | Show Queued only    |
| `i` | Show In Progress    |
| `s` | Show Succeeded      |
| `f` | Show Failed         |

### Search Modal
| Key          | Action             |
| ------------ | ------------------ |
| `Enter`      | Close search modal |
| `Escape`     | Close search modal |
| Clear Button | Clear search term  |
| Close Button | Close search modal |

> **Note**: Search filtering happens in real-time as you type and works in combination with status filters.

## Dashboard Layout

The terminal dashboard displays a comprehensive migration table with the following columns:

| Column          | Description                           |
| --------------- | ------------------------------------- |
| Repository Name | Name of the repository being migrated |
| Migration ID    | Unique identifier for the migration   |
| Status          | Current migration state (color-coded) |
| Created At      | When the migration was initiated      |

### Status Color Coding
- 🔵 **Blue**: Queued states (`QUEUED`, `WAITING`)
- 🟡 **Yellow**: In Progress (`IN_PROGRESS`, `PREPARING`, `PENDING`, `MAPPING`, `IMPORTING`, etc.)
- 🟢 **Green**: Succeeded (`SUCCEEDED`, `UNLOCKED`, `IMPORTED`)
- 🔴 **Red**: Failed (`FAILED`, `FAILED_IMPORT`)

### Smart Filtering & Search
- **Status Filters**: Instantly filter by migration state
- **Live Search**: Real-time repository name filtering
- **Combined Filtering**: Search works within selected status filters
- **Dynamic Title**: Shows organization name and active filter

## Performance & Optimization

### Recent Improvements
- **Optimized UI Code**: Refactored following KISS principles for better maintainability
- **Efficient Filtering**: Streamlined filter logic with single-responsibility functions
- **Memory Optimization**: Reduced object overhead and cleaner lifecycle management
- **Enhanced Search**: Real-time search with minimal performance impact
- **Faster Refresh**: Improved from 60-second to 30-second refresh intervals

### Technical Features
- **Separation of Concerns**: Clean architecture with focused, testable components
- **Event-Driven Updates**: Efficient handling of user interactions and data updates
- **Responsive Design**: Non-blocking UI updates and smooth animations
- **Error Resilience**: Graceful handling of API failures and network issues

## Requirements

- [GitHub CLI](https://cli.github.com) installed and authenticated
- GitHub token with permissions:
  - **GEI migrations**: `read:org`, `repo`
  - **Legacy migrations**: `read:org`, `repo`, `admin:org`

## User Interface Guide

### Getting Started
1. Launch with `gh migration-monitor --organization myorg`
2. The dashboard loads automatically with a 30-second refresh interval
3. Use keyboard shortcuts to navigate and filter data
4. Press `/` to search for specific repositories
5. Use status filters (`a`, `q`, `i`, `s`, `f`) to focus on specific migration states

### Search Functionality
- **Open**: Press `/` to open the search modal
- **Type**: Start typing a repository name for real-time filtering
- **Clear**: Click "Clear" button or manually delete text
- **Close**: Press `Enter`, `Escape`, or click "Close" button
- **Combine**: Search works with status filters for precise results

### Visual Indicators
- **Loading Animation**: Spinning indicator during data refresh
- **Last Updated**: Timestamp showing when data was last refreshed
- **Active Filter**: Current filter displayed in table title
- **Color Status**: Immediate visual status recognition

## Troubleshooting

**Token required**: Set `GHMM_GITHUB_TOKEN` or use `--github-token` flag

**Organization not found**: Check organization name and token permissions

**Rate limit**: Wait a few minutes - the tool automatically refreshes every 30 seconds

**Search not working**: Ensure you're in the main dashboard view (not in a modal)

**Keyboard shortcuts not responding**: Try pressing `Escape` to ensure you're not in search mode

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
3. Follow KISS principles for clean, maintainable code
4. Submit a pull request

### Project Structure
```
├── cmd/               # CLI commands (Cobra framework)
│   └── root.go       # Main command and application entry
├── internal/
│   ├── api/          # GitHub API clients (REST & GraphQL)
│   ├── config/       # Configuration management (Viper)
│   ├── models/       # Domain models and data structures
│   ├── services/     # Business logic and migration handling
│   └── ui/           # Terminal UI components (tview)
│       ├── ui.go     # Dashboard and interaction logic
│       └── table.go  # Migration table display
├── go.mod            # Go module definition
├── main.go           # Application entry point
└── README.md         # This documentation
```

### Architecture Highlights
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Interface-Based Design**: Testable and modular components
- **Real-time Updates**: Efficient background refresh with visual feedback
- **Responsive UI**: Optimized keyboard navigation and search functionality

## License

This project is licensed under the [MIT License](./LICENSE).

---

**Note**: This tool is inspired by [github-migration-monitor](https://github.com/timrogers/github-migration-monitor) and provides enhanced functionality with a modern terminal UI.
