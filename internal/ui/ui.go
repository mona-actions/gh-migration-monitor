package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/migration"
	"github.com/rivo/tview"
)

type UI struct {
	Queued     tview.Table
	InProgress tview.Table
	Succeeded  tview.Table
	Failed     tview.Table
}

func NewUI() *UI {
	return &UI{
		Queued:     *createTable("Queued"),
		InProgress: *createTable("In Progress"),
		Succeeded:  *createTable("Succeeded"),
		Failed:     *createTable("Failed"),
	}
}

func createTable(name string) *tview.Table {
	table := tview.NewTable().SetBorders(false)

	// Setting the table's box properties
	table.Box.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(name)

	return table
}

func (u *UI) UpdateData() {
	migrations := migration.GetMigrations()

	updateTable(&u.Queued, migrations.Queued)
	updateTable(&u.InProgress, migrations.In_Progress)
	updateTable(&u.Succeeded, migrations.Succeeded)
	updateTable(&u.Failed, migrations.Failed)
}

func updateTable(table *tview.Table, data []migration.Migration) {
	// Clear the table
	table.Clear()

	// Add the table's headers
	table.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
	table.SetCell(0, 1, tview.NewTableCell("Created At").SetExpansion(1))
	// Check for failed migration
	if table.GetTitle() == "Failed" {
		table.SetCell(0, 2, tview.NewTableCell("Error").SetExpansion(1))
	}

	// Add the table's data
	for i, migration := range data {
		table.SetCell(i+1, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(migration.CreatedAt).SetExpansion(1))
		// Check for failed migration
		if table.GetTitle() == "Failed" {
			table.SetCell(i+1, 2, tview.NewTableCell(migration.FailureReason).SetExpansion(1))
		}
	}
}
