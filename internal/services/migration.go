package services

import (
	"context"
	"fmt"

	"github.com/mona-actions/gh-migration-monitor/internal/api"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
)

// MigrationService handles migration-related business logic
type MigrationService interface {
	ListMigrations(ctx context.Context, org string, isLegacy bool) (*models.MigrationSummary, error)
}

// migrationService implements MigrationService
type migrationService struct {
	githubClient api.GitHubClient
}

// NewMigrationService creates a new migration service
func NewMigrationService(githubClient api.GitHubClient) MigrationService {
	return &migrationService{
		githubClient: githubClient,
	}
}

// ListMigrations retrieves and categorizes migrations by state
func (s *migrationService) ListMigrations(ctx context.Context, org string, isLegacy bool) (*models.MigrationSummary, error) {
	migrations, err := s.githubClient.ListMigrations(ctx, org, isLegacy)
	if err != nil {
		return nil, fmt.Errorf("failed to list migrations: %w", err)
	}

	summary := &models.MigrationSummary{
		Queued:     make([]models.Migration, 0),
		InProgress: make([]models.Migration, 0),
		Succeeded:  make([]models.Migration, 0),
		Failed:     make([]models.Migration, 0),
	}

	for _, migration := range migrations {
		switch {
		case migration.State.IsQueued():
			summary.Queued = append(summary.Queued, migration)
		case migration.State.IsInProgress():
			summary.InProgress = append(summary.InProgress, migration)
		case migration.State.IsSucceeded():
			summary.Succeeded = append(summary.Succeeded, migration)
		case migration.State.IsFailed():
			summary.Failed = append(summary.Failed, migration)
		}
	}

	return summary, nil
}
