package models

import "time"

// Migration represents a GitHub repository migration
type Migration struct {
	ID              string    `json:"id"`
	RepositoryName  string    `json:"repository_name"`
	State           State     `json:"state"`
	CreatedAt       time.Time `json:"created_at"`
	FailureReason   string    `json:"failure_reason,omitempty"`
	MigrationLogURL string    `json:"migration_log_url,omitempty"`
}

// State represents the current state of a migration
type State string

const (
	StateQueued       State = "QUEUED"
	StateWaiting      State = "WAITING"
	StateInProgress   State = "IN_PROGRESS"
	StatePreparing    State = "PREPARING"
	StatePending      State = "PENDING"
	StateMapping      State = "MAPPING"
	StateArchived     State = "ARCHIVE_UPLOADED"
	StateConflicts    State = "CONFLICTS"
	StateReady        State = "READY"
	StateImporting    State = "IMPORTING"
	StateSucceeded    State = "SUCCEEDED"
	StateUnlocked     State = "UNLOCKED"
	StateImported     State = "IMPORTED"
	StateFailed       State = "FAILED"
	StateFailedImport State = "FAILED_IMPORT"
)

// IsQueued returns true if the migration is in a queued state
func (s State) IsQueued() bool {
	return s == StateQueued || s == StateWaiting
}

// IsInProgress returns true if the migration is in progress
func (s State) IsInProgress() bool {
	return s == StateInProgress || s == StatePreparing || s == StatePending ||
		s == StateMapping || s == StateArchived || s == StateConflicts ||
		s == StateReady || s == StateImporting
}

// IsSucceeded returns true if the migration completed successfully
func (s State) IsSucceeded() bool {
	return s == StateSucceeded || s == StateUnlocked || s == StateImported
}

// IsFailed returns true if the migration failed
func (s State) IsFailed() bool {
	return s == StateFailed || s == StateFailedImport
}

// MigrationSummary provides a summary of migrations by state
type MigrationSummary struct {
	Queued     []Migration `json:"queued"`
	InProgress []Migration `json:"in_progress"`
	Succeeded  []Migration `json:"succeeded"`
	Failed     []Migration `json:"failed"`
}

// Total returns the total number of migrations
func (ms *MigrationSummary) Total() int {
	return len(ms.Queued) + len(ms.InProgress) + len(ms.Succeeded) + len(ms.Failed)
}

// ListOptions represents options for listing migrations
type ListOptions struct {
	Organization string `json:"organization"`
	State        State  `json:"state,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Page         int    `json:"page,omitempty"`
}
