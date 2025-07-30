# Contributing to GitHub Migration Monitor

Thank you for your interest in contributing to the GitHub Migration Monitor CLI extension! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.20 or later
- GitHub CLI installed and authenticated
- Git

### Development Setup

1. Fork and clone the repository:

   ```bash
   gh repo fork mona-actions/gh-migration-monitor --clone
   cd gh-migration-monitor
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Build the extension:

   ```bash
   go build -o gh-migration-monitor
   ```

4. Test the extension locally:
   ```bash
   gh extension install .
   gh migration-monitor --help
   ```

## Development Guidelines

### Code Organization

This project follows Go best practices and clean architecture principles:

- `cmd/`: CLI command definitions using Cobra
- `internal/api/`: GitHub API clients and integrations
- `internal/config/`: Configuration management with Viper
- `internal/models/`: Domain models and business entities
- `internal/services/`: Business logic services
- `internal/ui/`: Terminal UI components using tview

### Coding Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting
- Write tests for new functionality
- Document exported functions and types
- Use meaningful variable and function names
- Handle errors appropriately with context

### Testing

Run the test suite:

```bash
# Unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test ./internal/services/ -v
```

### Building and Testing Locally

```bash
# Build the binary
go build -o gh-migration-monitor

# Install as a local GitHub CLI extension
gh extension install .

# Test the extension
gh migration-monitor --organization your-test-org
```

## Submitting Changes

### Pull Request Process

1. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards

3. Add or update tests as needed

4. Run the test suite:

   ```bash
   go test ./...
   go vet ./...
   ```

5. Commit your changes with a clear message:

   ```bash
   git commit -m "Add feature: description of your changes"
   ```

6. Push to your fork:

   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a pull request with a clear title and description

### Pull Request Guidelines

- **Keep changes focused**: One feature or fix per pull request
- **Write clear commit messages**: Use the imperative mood (e.g., "Add feature" not "Added feature")
- **Update documentation**: Include relevant documentation updates
- **Add tests**: Ensure new functionality is tested
- **Follow Go conventions**: Use `gofmt`, `go vet`, and follow Go idioms

### Commit Message Format

Use clear, descriptive commit messages:

```
type: brief description (50 chars or less)

More detailed explanation if needed (wrap at 72 chars).
Explain what and why, not how.

- List any breaking changes
- Reference issues: Fixes #123
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

## Issues and Feature Requests

### Reporting Bugs

When reporting bugs, please include:

- Go version (`go version`)
- Operating system and version
- GitHub CLI version (`gh version`)
- Steps to reproduce the issue
- Expected vs actual behavior
- Error messages or logs

### Feature Requests

For feature requests, please describe:

- The problem you're trying to solve
- Your proposed solution
- Any alternative solutions considered
- Examples of how the feature would be used

## Development Tips

### Working with GitHub API

- Use the existing API clients in `internal/api/`
- Respect rate limits using the rate limiter
- Handle authentication through GitHub CLI when possible
- Use GraphQL for complex queries, REST for simple operations

### UI Development

- Follow tview patterns for terminal UI components
- Ensure accessibility and keyboard navigation
- Test UI components with different terminal sizes
- Use consistent styling and colors

### Configuration

- Use Viper for configuration management
- Support environment variables with `GH_MIGRATION_MONITOR_` prefix
- Provide sensible defaults
- Validate configuration on startup

## Code Review Process

1. All changes require code review
2. Maintainers will review pull requests within a few days
3. Address feedback promptly and respectfully
4. Once approved, maintainers will merge the changes

## Getting Help

- Check existing [issues](https://github.com/mona-actions/gh-migration-monitor/issues)
- Join discussions in pull requests
- Reach out to maintainers for guidance

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
