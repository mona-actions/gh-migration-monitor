package ui_test

import (
	"testing"
	"time"

	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/mona-actions/gh-migration-monitor/internal/ui"
	"github.com/rivo/tview"
)

func TestNewMigrationTable(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"queued table", "Queued"},
		{"in progress table", "In Progress"},
		{"succeeded table", "Succeeded"},
		{"failed table", "Failed"},
		{"empty title", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := ui.NewMigrationTable(tt.title)

			if table == nil {
				t.Fatal("NewMigrationTable() returned nil")
			}

			if table.Table == nil {
				t.Fatal("NewMigrationTable() returned table with nil Table field")
			}

			// Test that the table has the correct title
			if got := table.GetTitle(); got != tt.title {
				t.Errorf("NewMigrationTable() title = %v, want %v", got, tt.title)
			}

			// Test that the table is properly configured
			if table.Table == nil {
				t.Error("Table should not be nil")
			}
		})
	}
}

func TestMigrationTable_GetTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{"normal title", "Test Title", "Test Title"},
		{"empty title", "", ""},
		{"special characters", "Test-Title_123", "Test-Title_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := ui.NewMigrationTable(tt.title)
			if got := table.GetTitle(); got != tt.want {
				t.Errorf("MigrationTable.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigrationTable_UpdateData(t *testing.T) {
	now := time.Now()
	zeroTime := time.Time{}

	tests := []struct {
		name       string
		title      string
		migrations []models.Migration
		wantRows   int // Expected number of rows (including header)
		wantCols   int // Expected number of columns
	}{
		{
			name:       "empty migrations",
			title:      "Test",
			migrations: []models.Migration{},
			wantRows:   1, // Header only
			wantCols:   2, // Repository Name, Created At
		},
		{
			name:  "single migration with valid time",
			title: "Test",
			migrations: []models.Migration{
				{
					ID:             "1",
					RepositoryName: "test/repo1",
					State:          models.StateQueued,
					CreatedAt:      now,
				},
			},
			wantRows: 2, // Header + 1 data row
			wantCols: 2,
		},
		{
			name:  "multiple migrations",
			title: "Test",
			migrations: []models.Migration{
				{
					ID:             "1",
					RepositoryName: "test/repo1",
					State:          models.StateQueued,
					CreatedAt:      now,
				},
				{
					ID:             "2",
					RepositoryName: "test/repo2",
					State:          models.StateInProgress,
					CreatedAt:      zeroTime,
				},
			},
			wantRows: 3, // Header + 2 data rows
			wantCols: 2,
		},
		{
			name:  "failed migrations with error column",
			title: "Failed",
			migrations: []models.Migration{
				{
					ID:             "1",
					RepositoryName: "test/repo1",
					State:          models.StateFailed,
					CreatedAt:      now,
					FailureReason:  "Network timeout",
				},
				{
					ID:             "2",
					RepositoryName: "test/repo2",
					State:          models.StateFailed,
					CreatedAt:      now,
					FailureReason:  "",
				},
			},
			wantRows: 3, // Header + 2 data rows
			wantCols: 3, // Repository Name, Created At, Error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := ui.NewMigrationTable(tt.title)
			table.UpdateData(tt.migrations)

			// Check table dimensions
			rows := table.GetRowCount()
			cols := table.GetColumnCount()

			if rows != tt.wantRows {
				t.Errorf("UpdateData() rows = %v, want %v", rows, tt.wantRows)
			}

			if cols != tt.wantCols {
				t.Errorf("UpdateData() cols = %v, want %v", cols, tt.wantCols)
			}

			// Check header content
			if rows > 0 {
				repoCell := table.GetCell(0, 0)
				if repoCell.Text != "Repository Name" {
					t.Errorf("Header cell [0,0] = %v, want 'Repository Name'", repoCell.Text)
				}

				timeCell := table.GetCell(0, 1)
				if timeCell.Text != "Created At" {
					t.Errorf("Header cell [0,1] = %v, want 'Created At'", timeCell.Text)
				}

				// Check error column for failed migrations
				if tt.title == "Failed" && cols > 2 {
					errorCell := table.GetCell(0, 2)
					if errorCell.Text != "Error" {
						t.Errorf("Header cell [0,2] = %v, want 'Error'", errorCell.Text)
					}
				}
			}

			// Check data content
			for i, migration := range tt.migrations {
				row := i + 1
				if row >= rows {
					continue
				}

				// Check repository name
				repoCell := table.GetCell(row, 0)
				if repoCell.Text != migration.RepositoryName {
					t.Errorf("Data cell [%d,0] = %v, want %v", row, repoCell.Text, migration.RepositoryName)
				}

				// Check created at time formatting
				timeCell := table.GetCell(row, 1)
				expectedTime := "Unknown"
				if !migration.CreatedAt.IsZero() {
					expectedTime = migration.CreatedAt.Format("2006-01-02 15:04:05")
				}
				if timeCell.Text != expectedTime {
					t.Errorf("Data cell [%d,1] = %v, want %v", row, timeCell.Text, expectedTime)
				}

				// Check error column for failed migrations
				if tt.title == "Failed" && cols > 2 {
					errorCell := table.GetCell(row, 2)
					expectedError := migration.FailureReason
					if expectedError == "" {
						expectedError = "Unknown error"
					}
					if errorCell.Text != expectedError {
						t.Errorf("Data cell [%d,2] = %v, want %v", row, errorCell.Text, expectedError)
					}
				}
			}
		})
	}
}

func TestMigrationTable_UpdateData_ClearsExistingData(t *testing.T) {
	table := ui.NewMigrationTable("Test")

	// Add initial data
	initialMigrations := []models.Migration{
		{
			ID:             "1",
			RepositoryName: "test/repo1",
			State:          models.StateQueued,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "2",
			RepositoryName: "test/repo2",
			State:          models.StateQueued,
			CreatedAt:      time.Now(),
		},
	}
	table.UpdateData(initialMigrations)

	// Verify initial state
	if table.GetRowCount() != 3 { // Header + 2 rows
		t.Fatalf("Initial data not set correctly, got %d rows, want 3", table.GetRowCount())
	}

	// Update with new data
	newMigrations := []models.Migration{
		{
			ID:             "3",
			RepositoryName: "test/repo3",
			State:          models.StateQueued,
			CreatedAt:      time.Now(),
		},
	}
	table.UpdateData(newMigrations)

	// Verify old data is cleared and new data is set
	if table.GetRowCount() != 2 { // Header + 1 row
		t.Errorf("UpdateData() should clear existing data, got %d rows, want 2", table.GetRowCount())
	}

	// Verify the new data is correct
	if table.GetRowCount() > 1 {
		repoCell := table.GetCell(1, 0)
		if repoCell.Text != "test/repo3" {
			t.Errorf("New data not set correctly, got %v, want 'test/repo3'", repoCell.Text)
		}
	}
}

func TestMigrationTable_UpdateData_EmptyMigrations(t *testing.T) {
	table := ui.NewMigrationTable("Test")

	// Add some initial data
	initialMigrations := []models.Migration{
		{
			ID:             "1",
			RepositoryName: "test/repo1",
			State:          models.StateQueued,
			CreatedAt:      time.Now(),
		},
	}
	table.UpdateData(initialMigrations)

	// Update with empty migrations
	table.UpdateData([]models.Migration{})

	// Should only have header row
	if table.GetRowCount() != 1 {
		t.Errorf("UpdateData() with empty migrations should leave only header, got %d rows, want 1", table.GetRowCount())
	}

	// Verify header is still present
	repoCell := table.GetCell(0, 0)
	if repoCell.Text != "Repository Name" {
		t.Errorf("Header should still be present after clearing, got %v", repoCell.Text)
	}
}

// Test to ensure table implements the expected interface
func TestMigrationTable_ImplementsTableInterface(t *testing.T) {
	table := ui.NewMigrationTable("Test")

	// Verify it embeds *tview.Table
	if table.Table == nil {
		t.Error("MigrationTable should embed *tview.Table")
	}

	// Verify it can be used as a tview.Primitive
	var primitive tview.Primitive = table.Table
	_ = primitive // Verify assignment works
}

func TestMigrationTable_TimeFormatting(t *testing.T) {
	tests := []struct {
		name        string
		createdAt   time.Time
		expectedStr string
	}{
		{
			name:        "valid time",
			createdAt:   time.Date(2023, 7, 15, 14, 30, 45, 0, time.UTC),
			expectedStr: "2023-07-15 14:30:45",
		},
		{
			name:        "zero time",
			createdAt:   time.Time{},
			expectedStr: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := ui.NewMigrationTable("Test")
			migrations := []models.Migration{
				{
					ID:             "1",
					RepositoryName: "test/repo",
					State:          models.StateQueued,
					CreatedAt:      tt.createdAt,
				},
			}

			table.UpdateData(migrations)

			if table.GetRowCount() < 2 {
				t.Fatal("Table should have at least 2 rows (header + data)")
			}

			timeCell := table.GetCell(1, 1)
			if timeCell.Text != tt.expectedStr {
				t.Errorf("Time formatting = %v, want %v", timeCell.Text, tt.expectedStr)
			}
		})
	}
}
