// Package main provides the entry point for the gh-migration-monitor GitHub CLI extension.
//
// gh-migration-monitor is a terminal-based dashboard for monitoring GitHub organization
// migrations in real-time. It supports both legacy migrations and the new GitHub
// Enterprise Importer (GEI) migrations, providing an interactive interface to track
// migration progress across different states.
//
// The application follows clean architecture principles with clear separation of concerns:
//   - cmd/: Command-line interface and CLI parsing
//   - internal/api/: GitHub API client implementations
//   - internal/config/: Configuration management
//   - internal/models/: Domain models and business entities
//   - internal/services/: Business logic services
//   - internal/ui/: Terminal user interface components
//
// Usage:
//
//	go run main.go --organization myorg
//	go run main.go --organization myorg --legacy --github-token token
package main
