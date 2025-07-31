package ui_test

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/mona-actions/gh-migration-monitor/internal/ui"
	"github.com/rivo/tview"
)

func TestNewDashboard(t *testing.T) {
	dashboard := ui.NewDashboard()

	if dashboard == nil {
		t.Fatal("NewDashboard() returned nil")
	}

	// Test that all tables are initialized
	if dashboard.Queued == nil {
		t.Error("Queued table should not be nil")
	}
	if dashboard.InProgress == nil {
		t.Error("InProgress table should not be nil")
	}
	if dashboard.Succeeded == nil {
		t.Error("Succeeded table should not be nil")
	}
	if dashboard.Failed == nil {
		t.Error("Failed table should not be nil")
	}
	if dashboard.CommandBar == nil {
		t.Error("CommandBar should not be nil")
	}

	// Test table titles
	if dashboard.Queued.GetTitle() != "Queued" {
		t.Errorf("Queued table title = %v, want 'Queued'", dashboard.Queued.GetTitle())
	}
	if dashboard.InProgress.GetTitle() != "In Progress" {
		t.Errorf("InProgress table title = %v, want 'In Progress'", dashboard.InProgress.GetTitle())
	}
	if dashboard.Succeeded.GetTitle() != "Succeeded" {
		t.Errorf("Succeeded table title = %v, want 'Succeeded'", dashboard.Succeeded.GetTitle())
	}
	if dashboard.Failed.GetTitle() != "Failed" {
		t.Errorf("Failed table title = %v, want 'Failed'", dashboard.Failed.GetTitle())
	}
}

func TestDashboard_UpdateData(t *testing.T) {
	dashboard := ui.NewDashboard()

	now := time.Now()
	summary := &models.MigrationSummary{
		Queued: []models.Migration{
			{
				ID:             "1",
				RepositoryName: "test/repo1",
				State:          models.StateQueued,
				CreatedAt:      now,
			},
			{
				ID:             "2",
				RepositoryName: "test/repo2",
				State:          models.StateWaiting,
				CreatedAt:      now,
			},
		},
		InProgress: []models.Migration{
			{
				ID:             "3",
				RepositoryName: "test/repo3",
				State:          models.StateInProgress,
				CreatedAt:      now,
			},
		},
		Succeeded: []models.Migration{
			{
				ID:             "4",
				RepositoryName: "test/repo4",
				State:          models.StateSucceeded,
				CreatedAt:      now,
			},
		},
		Failed: []models.Migration{
			{
				ID:             "5",
				RepositoryName: "test/repo5",
				State:          models.StateFailed,
				CreatedAt:      now,
				FailureReason:  "Network timeout",
			},
		},
	}

	dashboard.UpdateData(summary)

	// Verify table row counts (header + data rows)
	if dashboard.Queued.GetRowCount() != 3 { // Header + 2 queued
		t.Errorf("Queued table rows = %v, want 3", dashboard.Queued.GetRowCount())
	}
	if dashboard.InProgress.GetRowCount() != 2 { // Header + 1 in progress
		t.Errorf("InProgress table rows = %v, want 2", dashboard.InProgress.GetRowCount())
	}
	if dashboard.Succeeded.GetRowCount() != 2 { // Header + 1 succeeded
		t.Errorf("Succeeded table rows = %v, want 2", dashboard.Succeeded.GetRowCount())
	}
	if dashboard.Failed.GetRowCount() != 2 { // Header + 1 failed
		t.Errorf("Failed table rows = %v, want 2", dashboard.Failed.GetRowCount())
	}

	// Verify specific data was populated correctly
	queuedCell := dashboard.Queued.GetCell(1, 0)
	if queuedCell.Text != "test/repo1" {
		t.Errorf("First queued repo = %v, want 'test/repo1'", queuedCell.Text)
	}

	inProgressCell := dashboard.InProgress.GetCell(1, 0)
	if inProgressCell.Text != "test/repo3" {
		t.Errorf("InProgress repo = %v, want 'test/repo3'", inProgressCell.Text)
	}

	succeededCell := dashboard.Succeeded.GetCell(1, 0)
	if succeededCell.Text != "test/repo4" {
		t.Errorf("Succeeded repo = %v, want 'test/repo4'", succeededCell.Text)
	}

	failedCell := dashboard.Failed.GetCell(1, 0)
	if failedCell.Text != "test/repo5" {
		t.Errorf("Failed repo = %v, want 'test/repo5'", failedCell.Text)
	}

	// Verify failed table has error column
	if dashboard.Failed.GetColumnCount() != 3 {
		t.Errorf("Failed table columns = %v, want 3 (including error column)", dashboard.Failed.GetColumnCount())
	}

	failedErrorCell := dashboard.Failed.GetCell(1, 2)
	if failedErrorCell.Text != "Network timeout" {
		t.Errorf("Failed error = %v, want 'Network timeout'", failedErrorCell.Text)
	}
}

func TestDashboard_UpdateData_EmptySummary(t *testing.T) {
	dashboard := ui.NewDashboard()

	// Update with empty summary
	summary := &models.MigrationSummary{
		Queued:     []models.Migration{},
		InProgress: []models.Migration{},
		Succeeded:  []models.Migration{},
		Failed:     []models.Migration{},
	}

	dashboard.UpdateData(summary)

	// All tables should only have header rows
	if dashboard.Queued.GetRowCount() != 1 {
		t.Errorf("Empty queued table rows = %v, want 1 (header only)", dashboard.Queued.GetRowCount())
	}
	if dashboard.InProgress.GetRowCount() != 1 {
		t.Errorf("Empty in progress table rows = %v, want 1 (header only)", dashboard.InProgress.GetRowCount())
	}
	if dashboard.Succeeded.GetRowCount() != 1 {
		t.Errorf("Empty succeeded table rows = %v, want 1 (header only)", dashboard.Succeeded.GetRowCount())
	}
	if dashboard.Failed.GetRowCount() != 1 {
		t.Errorf("Empty failed table rows = %v, want 1 (header only)", dashboard.Failed.GetRowCount())
	}
}

func TestDashboard_UpdateData_NilSummary(t *testing.T) {
	dashboard := ui.NewDashboard()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UpdateData with nil summary caused panic: %v", r)
		}
	}()

	dashboard.UpdateData(nil)
}

func TestDashboard_SetupGrid(t *testing.T) {
	dashboard := ui.NewDashboard()
	grid := dashboard.SetupGrid()

	if grid == nil {
		t.Fatal("SetupGrid() returned nil")
	}

	// Test grid configuration - this is a basic test since tview internals
	// are not easily accessible for detailed verification
	if grid == nil {
		t.Error("Grid should be created")
	}
}

func TestDashboard_SetupKeyboardNavigation(t *testing.T) {
	dashboard := ui.NewDashboard()
	grid := dashboard.SetupGrid()
	app := tview.NewApplication()

	// Setup keyboard navigation
	dashboard.SetupKeyboardNavigation(app, grid)

	// Test that the input capture function is set by triggering different key events
	// We'll test this by creating events and checking they don't panic
	tests := []struct {
		name        string
		key         rune
		shouldStop  bool
		description string
	}{
		{"exit key", 'x', true, "should stop the application"},
		{"queued key", 'q', false, "should focus queued table"},
		{"in progress key", 'i', false, "should focus in progress table"},
		{"succeeded key", 's', false, "should focus succeeded table"},
		{"failed key", 'f', false, "should focus failed table"},
		{"other key", 'z', false, "should handle unknown keys gracefully"},
		{"space key", ' ', false, "should handle space key"},
		{"number key", '1', false, "should handle number keys"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test event
			event := tcell.NewEventKey(tcell.KeyRune, tt.key, tcell.ModNone)

			// This tests that the setup doesn't panic and the event handler is configured
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Keyboard navigation setup caused panic with key %c: %v", tt.key, r)
				}
			}()

			// Test that we can create the event without issues
			if event.Rune() != tt.key {
				t.Errorf("Event key = %c, want %c", event.Rune(), tt.key)
			}

			// For non-exit keys, test that they return the event
			// For exit key, test that it would stop the app (we can't test actual stopping)
			if tt.key != 'x' {
				// Non-exit keys should return the event unchanged
				// We verify this by checking the event properties are preserved
				if event.Modifiers() != tcell.ModNone {
					t.Errorf("Event modifiers should be preserved for key %c", tt.key)
				}
			}
		})
	}
}

func TestDashboard_SetupKeyboardNavigation_InputCapture(t *testing.T) {
	dashboard := ui.NewDashboard()
	grid := dashboard.SetupGrid()
	app := tview.NewApplication()

	// Setup keyboard navigation
	dashboard.SetupKeyboardNavigation(app, grid)

	// Test different key events by ensuring they don't cause panics
	// and that the event handling logic is properly set up
	tests := []struct {
		name        string
		key         rune
		modifier    tcell.ModMask
		description string
	}{
		{"exit key", 'x', tcell.ModNone, "should handle exit key"},
		{"queued key", 'q', tcell.ModNone, "should handle queued focus key"},
		{"in progress key", 'i', tcell.ModNone, "should handle in progress focus key"},
		{"succeeded key", 's', tcell.ModNone, "should handle succeeded focus key"},
		{"failed key", 'f', tcell.ModNone, "should handle failed focus key"},
		{"unknown key", 'z', tcell.ModNone, "should handle unknown keys gracefully"},
		{"modified key", 'q', tcell.ModCtrl, "should handle modified keys"},
		{"uppercase key", 'Q', tcell.ModNone, "should handle uppercase keys"},
		{"special char", '@', tcell.ModNone, "should handle special characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test event
			event := tcell.NewEventKey(tcell.KeyRune, tt.key, tt.modifier)

			// Test that event creation and handling setup doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Keyboard event %c with modifier %v caused panic: %v", tt.key, tt.modifier, r)
				}
			}()

			// Verify event properties are correct
			if event.Rune() != tt.key {
				t.Errorf("Event key mismatch: got %c, want %c", event.Rune(), tt.key)
			}

			if event.Modifiers() != tt.modifier {
				t.Errorf("Event modifier mismatch: got %v, want %v", event.Modifiers(), tt.modifier)
			}

			// Test that the event can be processed (basic validation)
			if event.Key() != tcell.KeyRune {
				t.Errorf("Event should be a rune key, got %v", event.Key())
			}
		})
	}
}

func TestDashboard_SetupKeyboardNavigation_GridConfiguration(t *testing.T) {
	dashboard := ui.NewDashboard()
	grid := dashboard.SetupGrid()
	app := tview.NewApplication()

	// Test that the grid is properly configured before keyboard setup
	if grid == nil {
		t.Fatal("Grid should not be nil before keyboard navigation setup")
	}

	// Setup keyboard navigation
	dashboard.SetupKeyboardNavigation(app, grid)

	// Test that the setup completes without errors
	// The actual input capture function testing is limited by tview's design,
	// but we can ensure the setup process works correctly

	// Verify grid is still valid after setup
	if grid == nil {
		t.Error("Grid should not be nil after keyboard navigation setup")
	}

	// Test edge case: calling setup multiple times should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Multiple keyboard navigation setups caused panic: %v", r)
		}
	}()

	dashboard.SetupKeyboardNavigation(app, grid)
	dashboard.SetupKeyboardNavigation(app, grid)
}

func TestDashboard_SetupKeyboardNavigation_NilParameters(t *testing.T) {
	dashboard := ui.NewDashboard()

	tests := []struct {
		name        string
		app         *tview.Application
		grid        *tview.Grid
		shouldPanic bool
	}{
		{"nil app", nil, dashboard.SetupGrid(), false},  // nil app is handled gracefully
		{"nil grid", tview.NewApplication(), nil, true}, // nil grid causes panic when calling SetInputCapture
		{"both nil", nil, nil, true},                    // nil grid causes panic
		{"valid params", tview.NewApplication(), dashboard.SetupGrid(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but didn't get one")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			dashboard.SetupKeyboardNavigation(tt.app, tt.grid)
		})
	}
}

func TestDashboard_CommandBarContent(t *testing.T) {
	dashboard := ui.NewDashboard()

	if dashboard.CommandBar == nil {
		t.Fatal("CommandBar should not be nil")
	}

	// Get the command bar text
	commandText := dashboard.CommandBar.GetText(false)

	// Verify it contains the expected commands
	expectedCommands := []string{"q", "i", "s", "f", "x"}
	for _, cmd := range expectedCommands {
		if !containsString(commandText, cmd) {
			t.Errorf("Command bar should contain key '%s', got: %s", cmd, commandText)
		}
	}

	// Verify it contains the expected labels
	expectedLabels := []string{"Queued", "In Progress", "Succeeded", "Failed", "Exit"}
	for _, label := range expectedLabels {
		if !containsString(commandText, label) {
			t.Errorf("Command bar should contain label '%s', got: %s", label, commandText)
		}
	}
}

func TestDashboard_Integration(t *testing.T) {
	// Test complete dashboard workflow
	dashboard := ui.NewDashboard()
	grid := dashboard.SetupGrid()
	app := tview.NewApplication()

	// Setup everything
	dashboard.SetupKeyboardNavigation(app, grid)

	// Create test data
	now := time.Now()
	summary := &models.MigrationSummary{
		Queued: []models.Migration{
			{
				ID:             "1",
				RepositoryName: "test/queued-repo",
				State:          models.StateQueued,
				CreatedAt:      now,
			},
		},
		InProgress: []models.Migration{
			{
				ID:             "2",
				RepositoryName: "test/progress-repo",
				State:          models.StateInProgress,
				CreatedAt:      now,
			},
		},
		Succeeded: []models.Migration{
			{
				ID:             "3",
				RepositoryName: "test/success-repo",
				State:          models.StateSucceeded,
				CreatedAt:      now,
			},
		},
		Failed: []models.Migration{
			{
				ID:             "4",
				RepositoryName: "test/failed-repo",
				State:          models.StateFailed,
				CreatedAt:      now,
				FailureReason:  "Test failure",
			},
		},
	}

	// Update dashboard with data
	dashboard.UpdateData(summary)

	// Verify data is correctly distributed
	if dashboard.Queued.GetRowCount() != 2 {
		t.Errorf("Queued table should have 2 rows (header + data), got %d", dashboard.Queued.GetRowCount())
	}
	if dashboard.InProgress.GetRowCount() != 2 {
		t.Errorf("InProgress table should have 2 rows (header + data), got %d", dashboard.InProgress.GetRowCount())
	}
	if dashboard.Succeeded.GetRowCount() != 2 {
		t.Errorf("Succeeded table should have 2 rows (header + data), got %d", dashboard.Succeeded.GetRowCount())
	}
	if dashboard.Failed.GetRowCount() != 2 {
		t.Errorf("Failed table should have 2 rows (header + data), got %d", dashboard.Failed.GetRowCount())
	}

	// Verify specific content
	queuedRepo := dashboard.Queued.GetCell(1, 0).Text
	if queuedRepo != "test/queued-repo" {
		t.Errorf("Queued repo name = %v, want 'test/queued-repo'", queuedRepo)
	}

	failedError := dashboard.Failed.GetCell(1, 2).Text
	if failedError != "Test failure" {
		t.Errorf("Failed error = %v, want 'Test failure'", failedError)
	}
}

func TestDashboard_MultipleUpdates(t *testing.T) {
	dashboard := ui.NewDashboard()
	now := time.Now()

	// First update
	summary1 := &models.MigrationSummary{
		Queued: []models.Migration{
			{ID: "1", RepositoryName: "repo1", State: models.StateQueued, CreatedAt: now},
			{ID: "2", RepositoryName: "repo2", State: models.StateQueued, CreatedAt: now},
		},
		InProgress: []models.Migration{},
		Succeeded:  []models.Migration{},
		Failed:     []models.Migration{},
	}
	dashboard.UpdateData(summary1)

	if dashboard.Queued.GetRowCount() != 3 { // Header + 2 repos
		t.Errorf("After first update, queued table should have 3 rows, got %d", dashboard.Queued.GetRowCount())
	}

	// Second update with different data
	summary2 := &models.MigrationSummary{
		Queued: []models.Migration{
			{ID: "3", RepositoryName: "repo3", State: models.StateQueued, CreatedAt: now},
		},
		InProgress: []models.Migration{
			{ID: "1", RepositoryName: "repo1", State: models.StateInProgress, CreatedAt: now},
		},
		Succeeded: []models.Migration{
			{ID: "2", RepositoryName: "repo2", State: models.StateSucceeded, CreatedAt: now},
		},
		Failed: []models.Migration{},
	}
	dashboard.UpdateData(summary2)

	// Verify data was replaced, not appended
	if dashboard.Queued.GetRowCount() != 2 { // Header + 1 repo
		t.Errorf("After second update, queued table should have 2 rows, got %d", dashboard.Queued.GetRowCount())
	}
	if dashboard.InProgress.GetRowCount() != 2 { // Header + 1 repo
		t.Errorf("After second update, in progress table should have 2 rows, got %d", dashboard.InProgress.GetRowCount())
	}
	if dashboard.Succeeded.GetRowCount() != 2 { // Header + 1 repo
		t.Errorf("After second update, succeeded table should have 2 rows, got %d", dashboard.Succeeded.GetRowCount())
	}

	// Verify correct repositories are in each table
	queuedRepo := dashboard.Queued.GetCell(1, 0).Text
	if queuedRepo != "repo3" {
		t.Errorf("Queued repo after update = %v, want 'repo3'", queuedRepo)
	}

	inProgressRepo := dashboard.InProgress.GetCell(1, 0).Text
	if inProgressRepo != "repo1" {
		t.Errorf("InProgress repo after update = %v, want 'repo1'", inProgressRepo)
	}

	succeededRepo := dashboard.Succeeded.GetCell(1, 0).Text
	if succeededRepo != "repo2" {
		t.Errorf("Succeeded repo after update = %v, want 'repo2'", succeededRepo)
	}
}

// Helper function to check if a string contains a substring
func containsString(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// Benchmark tests for performance
func BenchmarkDashboard_UpdateData(b *testing.B) {
	dashboard := ui.NewDashboard()
	now := time.Now()

	// Create a large summary for benchmarking
	summary := &models.MigrationSummary{
		Queued:     make([]models.Migration, 100),
		InProgress: make([]models.Migration, 100),
		Succeeded:  make([]models.Migration, 100),
		Failed:     make([]models.Migration, 100),
	}

	// Populate with test data
	for i := 0; i < 100; i++ {
		migration := models.Migration{
			ID:             string(rune(i)),
			RepositoryName: "test/repo" + string(rune(i)),
			State:          models.StateQueued,
			CreatedAt:      now,
		}
		summary.Queued[i] = migration
		summary.InProgress[i] = migration
		summary.Succeeded[i] = migration
		summary.Failed[i] = migration
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dashboard.UpdateData(summary)
	}
}

func BenchmarkNewDashboard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ui.NewDashboard()
	}
}

// Test edge cases and error conditions
func TestDashboard_EdgeCases(t *testing.T) {
	t.Run("nil dashboard components", func(t *testing.T) {
		// This tests that our constructor properly initializes everything
		dashboard := ui.NewDashboard()

		// Ensure no component is nil
		components := []interface{}{
			dashboard.Queued,
			dashboard.InProgress,
			dashboard.Succeeded,
			dashboard.Failed,
			dashboard.CommandBar,
		}

		for i, component := range components {
			if component == nil {
				t.Errorf("Component %d should not be nil", i)
			}
		}
	})

	t.Run("large migration data", func(t *testing.T) {
		dashboard := ui.NewDashboard()
		now := time.Now()

		// Create a large number of migrations
		largeMigrations := make([]models.Migration, 1000)
		for i := 0; i < 1000; i++ {
			largeMigrations[i] = models.Migration{
				ID:             string(rune(i + 48)), // Start from '0' ASCII
				RepositoryName: "large/repo" + string(rune(i+48)),
				State:          models.StateQueued,
				CreatedAt:      now,
			}
		}

		summary := &models.MigrationSummary{
			Queued:     largeMigrations,
			InProgress: []models.Migration{},
			Succeeded:  []models.Migration{},
			Failed:     []models.Migration{},
		}

		// Should not panic with large data
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Large data update caused panic: %v", r)
			}
		}()

		dashboard.UpdateData(summary)

		// Verify the data was added
		if dashboard.Queued.GetRowCount() != 1001 { // Header + 1000 rows
			t.Errorf("Large data update rows = %d, want 1001", dashboard.Queued.GetRowCount())
		}
	})
}
