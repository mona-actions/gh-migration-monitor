package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/rivo/tview"
)

// FilterOption represents different filter options
type FilterOption string

const (
	FilterAll        FilterOption = "All"
	FilterQueued     FilterOption = "Queued"
	FilterInProgress FilterOption = "In Progress"
	FilterSucceeded  FilterOption = "Succeeded"
	FilterFailed     FilterOption = "Failed"
)

// Dashboard represents the main UI dashboard
type Dashboard struct {
	AllMigrations    *MigrationTable
	CommandBar       *tview.TextView
	StatusBar        *tview.TextView
	SearchInput      *tview.InputField
	MainGrid         *tview.Grid
	app              *tview.Application
	refreshFunc      func()
	isRefreshing     bool
	currentFilter    FilterOption
	allMigrations    []models.Migration
	organizationName string
	searchTerm       string
}

// NewDashboard creates a new UI dashboard
func NewDashboard() *Dashboard {
	dashboard := &Dashboard{
		AllMigrations: NewMigrationTable("Migration Status"),
		CommandBar:    createCommandBar(),
		StatusBar:     createStatusBar(),
		isRefreshing:  false,
		currentFilter: FilterAll,
		allMigrations: make([]models.Migration, 0),
		searchTerm:    "",
	}

	// Create search input
	dashboard.SearchInput = createSearchInput()

	return dashboard
}

// createCommandBar creates a text view displaying keyboard shortcuts
func createCommandBar() *tview.TextView {
	commandBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow::b]Commands: [white::]r[grey::] Refresh  [white::]/ [grey::] Search  [white::]x[grey::] Exit  [yellow::b]Filters: [white::]a[grey::] All  [white::]q[grey::] Queued  [white::]i[grey::] In Progress  [white::]s[grey::] Succeeded  [white::]f[grey::] Failed")

	commandBar.SetBorder(false)

	return commandBar
}

// createSearchInput creates the search input field
func createSearchInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Search: ").
		SetPlaceholder("Type repository name...").
		SetFieldWidth(30)
}

// createStatusBar creates a text view for displaying status and progress
func createStatusBar() *tview.TextView {
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetText("")

	statusBar.SetBorder(false)

	return statusBar
}

// UpdateData updates the table with new migration data
func (d *Dashboard) UpdateData(summary *models.MigrationSummary, organization string) {
	if summary == nil {
		return
	}

	// Store organization name
	d.organizationName = organization

	// Update table title with organization name and current filter
	d.updateTitle()

	// Combine all migrations into a single list and store them
	d.allMigrations = make([]models.Migration, 0)
	d.allMigrations = append(d.allMigrations, summary.Queued...)
	d.allMigrations = append(d.allMigrations, summary.InProgress...)
	d.allMigrations = append(d.allMigrations, summary.Succeeded...)
	d.allMigrations = append(d.allMigrations, summary.Failed...)

	// Apply current filter
	d.applyFilter()
}

// applyFilter filters the migrations based on the current filter setting and search term
func (d *Dashboard) applyFilter() {
	if len(d.allMigrations) == 0 {
		d.AllMigrations.UpdateDataWithStatus([]models.Migration{})
		return
	}

	filteredMigrations := d.filterByStatus()
	filteredMigrations = d.filterBySearch(filteredMigrations)

	d.AllMigrations.UpdateDataWithStatus(filteredMigrations)
}

// filterByStatus filters migrations by status
func (d *Dashboard) filterByStatus() []models.Migration {
	if d.currentFilter == FilterAll {
		return d.allMigrations
	}

	var filtered []models.Migration
	for _, migration := range d.allMigrations {
		if d.matchesCurrentFilter(migration) {
			filtered = append(filtered, migration)
		}
	}
	return filtered
}

// filterBySearch filters migrations by search term
func (d *Dashboard) filterBySearch(migrations []models.Migration) []models.Migration {
	if d.searchTerm == "" {
		return migrations
	}

	var filtered []models.Migration
	searchLower := strings.ToLower(d.searchTerm)

	for _, migration := range migrations {
		if strings.Contains(strings.ToLower(migration.RepositoryName), searchLower) {
			filtered = append(filtered, migration)
		}
	}
	return filtered
}

// matchesCurrentFilter checks if a migration matches the current status filter
func (d *Dashboard) matchesCurrentFilter(migration models.Migration) bool {
	switch d.currentFilter {
	case FilterQueued:
		return migration.State.IsQueued()
	case FilterInProgress:
		return migration.State.IsInProgress()
	case FilterSucceeded:
		return migration.State.IsSucceeded()
	case FilterFailed:
		return migration.State.IsFailed()
	default:
		return true
	}
}

// SetupGrid creates and configures the grid layout
func (d *Dashboard) SetupGrid() *tview.Grid {
	if d.MainGrid == nil {
		d.MainGrid = tview.NewGrid().
			SetRows(0, 1).
			SetColumns(0).
			SetBorders(false)

		// Add the main migration table
		d.MainGrid.AddItem(d.AllMigrations.Table, 0, 0, 1, 1, 0, 0, true)

		// Create a flex layout for the bottom row containing command bar and status bar
		bottomFlex := tview.NewFlex().
			AddItem(d.CommandBar, 0, 3, false).
			AddItem(d.StatusBar, 0, 1, false)

		// Add bottom flex at the bottom with fixed height of 1 row
		d.MainGrid.AddItem(bottomFlex, 1, 0, 1, 1, 0, 0, false)
	}

	return d.MainGrid
}

// SetupKeyboardNavigation configures keyboard event handling
func (d *Dashboard) SetupKeyboardNavigation(app *tview.Application, grid *tview.Grid) {
	d.app = app
	grid.SetInputCapture(d.handleKeyInput)
}

// handleKeyInput processes keyboard input for the main dashboard
func (d *Dashboard) handleKeyInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'x':
		d.app.Stop()
		return nil
	case 'r':
		d.handleRefresh()
	case '/':
		d.showSearchModal()
		return nil
	case 'a', 'q', 'i', 's', 'f':
		d.handleFilterKey(event.Rune())
		return nil
	}
	return event
}

// handleRefresh triggers a refresh if one is not already in progress
func (d *Dashboard) handleRefresh() {
	if d.refreshFunc != nil && !d.isRefreshing {
		go d.refreshFunc()
	}
}

// handleFilterKey processes filter shortcut keys
func (d *Dashboard) handleFilterKey(key rune) {
	filterMap := map[rune]FilterOption{
		'a': FilterAll,
		'q': FilterQueued,
		'i': FilterInProgress,
		's': FilterSucceeded,
		'f': FilterFailed,
	}

	if filter, exists := filterMap[key]; exists {
		d.setFilter(filter)
	}
}

// setFilter sets the current filter and updates the display
func (d *Dashboard) setFilter(filter FilterOption) {
	d.currentFilter = filter

	// Update title and apply the filter
	d.updateTitle()
	d.applyFilter()
}

// updateTitle updates the table title with organization and current filter
func (d *Dashboard) updateTitle() {
	if d.organizationName != "" {
		d.AllMigrations.SetTitleWithOrganizationAndFilter(d.organizationName, string(d.currentFilter))
	}
}

// showSearchModal displays the search modal
func (d *Dashboard) showSearchModal() {
	if d.app == nil {
		return
	}

	// Setup search input
	d.SearchInput.SetText(d.searchTerm)
	d.setupSearchHandlers()

	// Create and show search form
	searchForm := d.createSearchForm()
	pages := tview.NewPages().
		AddPage("main", d.MainGrid, true, true).
		AddPage("search", searchForm, true, true)

	d.app.SetRoot(pages, true)
	d.app.SetFocus(d.SearchInput)
}

// setupSearchHandlers configures the search input event handlers
func (d *Dashboard) setupSearchHandlers() {
	d.SearchInput.SetChangedFunc(func(text string) {
		d.searchTerm = text
		d.applyFilter()
	})

	d.SearchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter || key == tcell.KeyEscape {
			d.app.QueueUpdateDraw(d.closeSearchModal)
		}
	})
}

// createSearchForm creates the search form layout
func (d *Dashboard) createSearchForm() tview.Primitive {
	form := tview.NewForm().
		AddFormItem(d.SearchInput).
		AddButton("Clear", d.clearSearch).
		AddButton("Close", d.closeSearchModal)

	form.SetBorder(true).
		SetTitle(" Search Repositories ").
		SetBorderColor(tcell.ColorTeal).
		SetInputCapture(d.handleFormInput)

	return d.centerForm(form)
}

// clearSearch clears the search term and updates the display
func (d *Dashboard) clearSearch() {
	d.SearchInput.SetText("")
	d.searchTerm = ""
	d.applyFilter()
}

// handleFormInput handles keyboard input for the search form
func (d *Dashboard) handleFormInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEscape {
		d.closeSearchModal()
		return nil
	}
	return event
}

// centerForm centers the form in the available space
func (d *Dashboard) centerForm(form *tview.Form) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 7, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)
}

// closeSearchModal hides the search modal and returns to main view
func (d *Dashboard) closeSearchModal() {
	if d.app == nil || d.MainGrid == nil {
		return
	}

	// Clear search input handlers
	d.SearchInput.SetChangedFunc(nil)
	d.SearchInput.SetDoneFunc(nil)

	// Restore main view
	d.app.SetRoot(d.MainGrid, true)
	d.SetupKeyboardNavigation(d.app, d.MainGrid)
	d.app.SetFocus(d.AllMigrations.Table)
}

// SetRefreshFunc sets the function to call when refresh is triggered
func (d *Dashboard) SetRefreshFunc(f func()) {
	d.refreshFunc = f
}

// ShowRefreshing displays a loading indicator with animation
func (d *Dashboard) ShowRefreshing() {
	if d.app != nil {
		d.isRefreshing = true

		// Start animation
		go func() {
			frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			frameIndex := 0

			for d.isRefreshing {
				d.app.QueueUpdateDraw(func() {
					if d.isRefreshing {
						d.StatusBar.SetText(fmt.Sprintf("[yellow::b]%s Refreshing...", frames[frameIndex]))
					}
				})
				frameIndex = (frameIndex + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
}

// HideRefreshing hides the loading indicator and shows last update time
func (d *Dashboard) HideRefreshing() {
	if d.app != nil {
		d.isRefreshing = false
		d.app.QueueUpdateDraw(func() {
			currentTime := time.Now().Format("15:04:05")
			d.StatusBar.SetText(fmt.Sprintf("[green::b]Last updated: %s", currentTime))
		})
	}
}

// ShowProgress shows progress with animated dots
func (d *Dashboard) ShowProgress(message string) {
	if d.app != nil {
		d.app.QueueUpdateDraw(func() {
			d.StatusBar.SetText(fmt.Sprintf("[yellow::b]%s", message))
		})
	}
}
