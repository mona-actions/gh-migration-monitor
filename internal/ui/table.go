package ui

import (
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

	// Add error column for failed migrations
	if mt.title == "Failed" {
		mt.SetCell(0, 2, tview.NewTableCell("Error").SetExpansion(1))
	}

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

		// Add error information for failed migrations
		if mt.title == "Failed" {
			errorText := migration.FailureReason
			if errorText == "" {
				errorText = "Unknown error"
			}
			mt.SetCell(row, 2, tview.NewTableCell(errorText).SetExpansion(1))
		}
	}
}

// GetTitle returns the table title
func (mt *MigrationTable) GetTitle() string {
	return mt.title
}
