package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/rivo/tview"
)

// MigrationTable represents a table for displaying migrations
type MigrationTable struct {
	*tview.Table
	title string
}

// NewMigrationTable creates a new migration table
func NewMigrationTable(title string) *MigrationTable {
	table := tview.NewTable().SetBorders(false)
	table.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)

	return &MigrationTable{
		Table: table,
		title: title,
	}
}

// UpdateData updates the table with new migration data
func (mt *MigrationTable) UpdateData(migrations []models.Migration) {
	mt.Clear()

	// Add headers
	mt.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
	mt.SetCell(0, 1, tview.NewTableCell("Created At").SetExpansion(1))

	// Add migration data
	for i, migration := range migrations {
		row := i + 1
		mt.SetCell(row, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))

		// Format the created at time
		formattedTime := migration.CreatedAt.Format("2006-01-02 15:04:05")
		if migration.CreatedAt.IsZero() {
			formattedTime = "Unknown"
		}
		mt.SetCell(row, 1, tview.NewTableCell(formattedTime).SetExpansion(1))
	}
}

// UpdateDataWithStatus updates the table with migration data including status information
func (mt *MigrationTable) UpdateDataWithStatus(migrations []models.Migration) {
	mt.Clear()

	// Add headers
	mt.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
	mt.SetCell(0, 1, tview.NewTableCell("Migration ID").SetExpansion(1))
	mt.SetCell(0, 2, tview.NewTableCell("Status").SetExpansion(1))
	mt.SetCell(0, 3, tview.NewTableCell("Created At").SetExpansion(1))

	// Add migration data
	for i, migration := range migrations {
		row := i + 1

		// Repository Name column
		mt.SetCell(row, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))

		// Migration ID column
		mt.SetCell(row, 1, tview.NewTableCell(migration.ID).SetExpansion(1))

		// Add status with color coding
		status := string(migration.State)
		statusCell := tview.NewTableCell(status).SetExpansion(1)

		// Color code the status
		switch {
		case migration.State.IsSucceeded():
			statusCell.SetTextColor(tcell.ColorGreen)
		case migration.State.IsFailed():
			statusCell.SetTextColor(tcell.ColorRed)
		case migration.State.IsInProgress():
			statusCell.SetTextColor(tcell.ColorYellow)
		case migration.State.IsQueued():
			statusCell.SetTextColor(tcell.ColorBlue)
		default:
			statusCell.SetTextColor(tcell.ColorWhite)
		}
		mt.SetCell(row, 2, statusCell)

		// Format the created at time
		formattedTime := migration.CreatedAt.Format("2006-01-02 15:04:05")
		if migration.CreatedAt.IsZero() {
			formattedTime = "Unknown"
		}
		mt.SetCell(row, 3, tview.NewTableCell(formattedTime).SetExpansion(1))
	}
}

// GetTitle returns the table title
func (mt *MigrationTable) GetTitle() string {
	return mt.title
}

// SetTitleWithOrganizationAndFilter updates the table title to include organization and filter
func (mt *MigrationTable) SetTitleWithOrganizationAndFilter(organization, filter string) {
	var newTitle string
	if filter == "All" || filter == "" {
		newTitle = fmt.Sprintf("Migration Status - %s", organization)
	} else {
		newTitle = fmt.Sprintf("Migration Status - %s [%s]", organization, filter)
	}
	mt.title = newTitle
	mt.Table.SetTitle(newTitle)
}
