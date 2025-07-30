package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/rivo/tview"
)

// Dashboard represents the main UI dashboard
type Dashboard struct {
	Queued     *MigrationTable
	InProgress *MigrationTable
	Succeeded  *MigrationTable
	Failed     *MigrationTable
	CommandBar *tview.TextView
}

// NewDashboard creates a new UI dashboard
func NewDashboard() *Dashboard {
	return &Dashboard{
		Queued:     NewMigrationTable("Queued"),
		InProgress: NewMigrationTable("In Progress"),
		Succeeded:  NewMigrationTable("Succeeded"),
		Failed:     NewMigrationTable("Failed"),
		CommandBar: createCommandBar(),
	}
}

// createCommandBar creates a text view displaying keyboard shortcuts
func createCommandBar() *tview.TextView {
	commandBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow::b]Commands: [white::]q[grey::] Queued  [white::]i[grey::] In Progress  [white::]s[grey::] Succeeded  [white::]f[grey::] Failed  [white::]x[grey::] Exit")

	commandBar.SetBorder(false)

	return commandBar
}

// UpdateData updates all tables with new migration data
func (d *Dashboard) UpdateData(summary *models.MigrationSummary) {
	if summary == nil {
		return
	}
	d.Queued.UpdateData(summary.Queued)
	d.InProgress.UpdateData(summary.InProgress)
	d.Succeeded.UpdateData(summary.Succeeded)
	d.Failed.UpdateData(summary.Failed)
}

// SetupGrid creates and configures the grid layout
func (d *Dashboard) SetupGrid() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(0, 0, 1).
		SetColumns(0, 0, 0).
		SetBorders(false)

	// Add tables to the grid
	grid.AddItem(d.Queued.Table, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(d.InProgress.Table, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(d.Succeeded.Table, 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(d.Failed.Table, 1, 0, 1, 3, 0, 0, false)

	// Add command bar at the bottom with fixed height of 1 row
	grid.AddItem(d.CommandBar, 2, 0, 1, 3, 0, 0, false)

	return grid
}

// SetupKeyboardNavigation configures keyboard event handling
func (d *Dashboard) SetupKeyboardNavigation(app *tview.Application, grid *tview.Grid) {
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'x':
			app.Stop()
			return nil
		case 'q':
			app.SetFocus(d.Queued.Table)
		case 'i':
			app.SetFocus(d.InProgress.Table)
		case 's':
			app.SetFocus(d.Succeeded.Table)
		case 'f':
			app.SetFocus(d.Failed.Table)
		}
		return event
	})
}
