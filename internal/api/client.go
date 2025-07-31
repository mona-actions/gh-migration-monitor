package api

import (
	"context"

	"github.com/mona-actions/gh-migration-monitor/internal/models"
)

// GitHubClient defines the interface for GitHub API operations
type GitHubClient interface {
	// ListMigrations returns all migrations for the specified organization
	ListMigrations(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error)
}

// APIError represents an API error
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}
