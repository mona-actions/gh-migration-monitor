package models_test

import (
	"testing"

	"github.com/mona-actions/gh-migration-monitor/internal/models"
)

func TestState_IsQueued(t *testing.T) {
	tests := []struct {
		name  string
		state models.State
		want  bool
	}{
		{"QUEUED", models.StateQueued, true},
		{"WAITING", models.StateWaiting, true},
		{"IN_PROGRESS", models.StateInProgress, false},
		{"SUCCEEDED", models.StateSucceeded, false},
		{"FAILED", models.StateFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsQueued(); got != tt.want {
				t.Errorf("State.IsQueued() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_IsInProgress(t *testing.T) {
	tests := []struct {
		name  string
		state models.State
		want  bool
	}{
		{"IN_PROGRESS", models.StateInProgress, true},
		{"PREPARING", models.StatePreparing, true},
		{"IMPORTING", models.StateImporting, true},
		{"QUEUED", models.StateQueued, false},
		{"SUCCEEDED", models.StateSucceeded, false},
		{"FAILED", models.StateFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsInProgress(); got != tt.want {
				t.Errorf("State.IsInProgress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_IsSucceeded(t *testing.T) {
	tests := []struct {
		name  string
		state models.State
		want  bool
	}{
		{"SUCCEEDED", models.StateSucceeded, true},
		{"UNLOCKED", models.StateUnlocked, true},
		{"IMPORTED", models.StateImported, true},
		{"QUEUED", models.StateQueued, false},
		{"IN_PROGRESS", models.StateInProgress, false},
		{"FAILED", models.StateFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsSucceeded(); got != tt.want {
				t.Errorf("State.IsSucceeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_IsFailed(t *testing.T) {
	tests := []struct {
		name  string
		state models.State
		want  bool
	}{
		{"FAILED", models.StateFailed, true},
		{"FAILED_IMPORT", models.StateFailedImport, true},
		{"QUEUED", models.StateQueued, false},
		{"IN_PROGRESS", models.StateInProgress, false},
		{"SUCCEEDED", models.StateSucceeded, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsFailed(); got != tt.want {
				t.Errorf("State.IsFailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigrationSummary_Total(t *testing.T) {
	summary := &models.MigrationSummary{
		Queued:     make([]models.Migration, 2),
		InProgress: make([]models.Migration, 3),
		Succeeded:  make([]models.Migration, 5),
		Failed:     make([]models.Migration, 1),
	}

	want := 11
	if got := summary.Total(); got != want {
		t.Errorf("MigrationSummary.Total() = %v, want %v", got, want)
	}
}
