package monitor

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/ui"
	"github.com/rivo/tview"
)

func Organization() {
	// Initialize the application
	app := tview.NewApplication()

	// Create the main layout
	grid := tview.NewGrid().
		SetSize(1, 2, 0, 0).
		SetBorders(false)

	// Adding tables to the grid
	tables := ui.NewUI()
	grid.AddItem(&tables.Queued, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(&tables.InProgress, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(&tables.Succeeded, 0, 2, 1, 1, 0, 0, false)

	// Capture keyboard events
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Quit app if user presses 'q'
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	// Update the data every 5 seconds
	go func() {
		// Testing the UI
		app.QueueUpdateDraw(func() {
			fmt.Println("Updating data...", time.Now())
		})
	}()

	// Set the grid as the application root
	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
