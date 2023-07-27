package monitor

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/migration"
	"github.com/rivo/tview"
)

var (
	app        *tview.Application
	queued     *tview.Table
	inProgress *tview.Table
	succeeded  *tview.Table
	logs       *tview.Table
)

func updateMigrations(interval int) {

	// Fetching the data
	var migrations migration.Migrations
	migrations.FetchMigrations()

	// Populating the tables
	// Yes, I know this is not the best way to do it, I hate this code every time I look at it
	app.QueueUpdateDraw(func() {
		queued.Clear()
		queued.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
		queued.SetCell(0, 1, tview.NewTableCell("Created At").SetExpansion(1))

		inProgress.Clear()
		inProgress.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
		inProgress.SetCell(0, 1, tview.NewTableCell("Created At").SetExpansion(1))

		succeeded.Clear()
		succeeded.SetCell(0, 0, tview.NewTableCell("Repository Name").SetExpansion(1))
		succeeded.SetCell(0, 1, tview.NewTableCell("Created At").SetExpansion(1))

		logs.SetCell(0, 0, tview.NewTableCell("").SetExpansion(1))

		for i, migration := range migrations.Queued {
			queued.SetCell(i+1, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))
			queued.SetCell(i+1, 1, tview.NewTableCell(migration.CreatedAt).SetExpansion(1))
		}

		for i, migration := range migrations.In_Progress {
			inProgress.SetCell(i+1, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))
			inProgress.SetCell(i+1, 1, tview.NewTableCell(migration.CreatedAt).SetExpansion(1))
		}

		for i, migration := range migrations.Succeeded {
			succeeded.SetCell(i+1, 0, tview.NewTableCell(migration.RepositoryName).SetExpansion(1))
			succeeded.SetCell(i+1, 1, tview.NewTableCell(migration.CreatedAt).SetExpansion(1))
		}

		// logs.SetCell(logs.GetRowCount()-1, 0, tview.NewTableCell(errors.Error()).SetExpansion(1))
	})

	// Refresh Interval
	time.Sleep(time.Duration(interval) * time.Second)
}

func Organization() {

	// Initialize the application
	app = tview.NewApplication()

	newBoxPrimitive := func(text string) tview.Primitive {
		return tview.NewBox().
			SetBorder(true).
			SetBorderColor(tcell.ColorTeal).
			SetTitleAlign(tview.AlignLeft).
			SetTitle(text)
	}

	grid := tview.NewGrid().
		SetSize(3, 3, 0, 0).
		SetBorders(false)

	// Building Table
	queued = tview.NewTable().SetBorders(false)
	queued.Box.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Queued")

	inProgress = tview.NewTable().SetBorders(false)
	inProgress.Box.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("In Progress")

	succeeded = tview.NewTable().SetBorders(false)
	succeeded.Box.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Succeeded")

	// TODO: Add Failed table

	logs = tview.NewTable().SetBorders(false)
	logs.Box.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Logs")

	// Add Status Boxes
	grid.AddItem(queued, 0, 0, 1, 1, 0, 100, false).
		AddItem(inProgress, 0, 1, 1, 1, 0, 100, false).
		AddItem(succeeded, 0, 2, 1, 1, 0, 100, false)

	// Add Failed Box
	grid.AddItem(newBoxPrimitive("Failed"), 1, 0, 1, 3, 0, 100, false)

	// Add Log Box
	grid.AddItem(logs, 2, 0, 1, 3, 0, 100, false)

	// Capture user input
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Quit app if user presses 'q'
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	// Update migrations every 5 seconds
	go updateMigrations(5)

	// Set the grid as the application root
	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
