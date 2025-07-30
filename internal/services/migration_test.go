package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/mona-actions/gh-migration-monitor/internal/services"
)

// MockGitHubClient implements the GitHubClient interface for testing
type MockGitHubClient struct {
	ListMigrationsFunc func(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error)
}

func (m *MockGitHubClient) ListMigrations(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error) {
	if m.ListMigrationsFunc != nil {
		return m.ListMigrationsFunc(ctx, org, isLegacy)
	}
	return []models.Migration{}, nil
}

func TestMigrationService_ListMigrations(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*MockGitHubClient)
		org        string
		isLegacy   bool
		wantQueued int
		wantInProg int
		wantSucc   int
		wantFailed int
		wantErr    bool
	}{
		{
			name: "successful categorization",
			setupMock: func(m *MockGitHubClient) {
				m.ListMigrationsFunc = func(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error) {
					return []models.Migration{
						{ID: "1", State: models.StateQueued, RepositoryName: "repo1", CreatedAt: time.Now()},
						{ID: "2", State: models.StateInProgress, RepositoryName: "repo2", CreatedAt: time.Now()},
						{ID: "3", State: models.StateSucceeded, RepositoryName: "repo3", CreatedAt: time.Now()},
						{ID: "4", State: models.StateFailed, RepositoryName: "repo4", CreatedAt: time.Now()},
					}, nil
				}
			},
			org:        "testorg",
			isLegacy:   false,
			wantQueued: 1,
			wantInProg: 1,
			wantSucc:   1,
			wantFailed: 1,
			wantErr:    false,
		},
		{
			name: "empty result",
			setupMock: func(m *MockGitHubClient) {
				m.ListMigrationsFunc = func(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error) {
					return []models.Migration{}, nil
				}
			},
			org:        "testorg",
			isLegacy:   false,
			wantQueued: 0,
			wantInProg: 0,
			wantSucc:   0,
			wantFailed: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockGitHubClient{}
			tt.setupMock(mockClient)

			service := services.NewMigrationService(mockClient)
			got, err := service.ListMigrations(context.Background(), tt.org, tt.isLegacy)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListMigrations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if len(got.Queued) != tt.wantQueued {
				t.Errorf("ListMigrations() queued = %v, want %v", len(got.Queued), tt.wantQueued)
			}

			if len(got.InProgress) != tt.wantInProg {
				t.Errorf("ListMigrations() in progress = %v, want %v", len(got.InProgress), tt.wantInProg)
			}

			if len(got.Succeeded) != tt.wantSucc {
				t.Errorf("ListMigrations() succeeded = %v, want %v", len(got.Succeeded), tt.wantSucc)
			}

			if len(got.Failed) != tt.wantFailed {
				t.Errorf("ListMigrations() failed = %v, want %v", len(got.Failed), tt.wantFailed)
			}
		})
	}
}
