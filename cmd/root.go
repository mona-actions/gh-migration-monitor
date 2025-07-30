package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mona-actions/gh-migration-monitor/internal/api"
	"github.com/mona-actions/gh-migration-monitor/internal/config"
	"github.com/mona-actions/gh-migration-monitor/internal/services"
	"github.com/mona-actions/gh-migration-monitor/internal/ui"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var (
	organization string
	githubToken  string
	legacy       bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "migration-monitor",
	Short: "GitHub CLI extension to monitor migration status",
	Long: `A GitHub CLI extension that monitors the progress of GitHub Organization migrations.

This tool provides a real-time dashboard for tracking repository migrations, supporting both
legacy migrations and the new GitHub Enterprise Importer (GEI) migrations.`,
	RunE: runMigrationMonitor,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Required flags
	rootCmd.Flags().StringVarP(&organization, "organization", "o", "", "GitHub organization to monitor (required)")

	// Optional flags
	rootCmd.Flags().StringVarP(&githubToken, "github-token", "t", "", "GitHub token (can also be set via GHMM_GITHUB_TOKEN)")
	rootCmd.Flags().BoolVarP(&legacy, "legacy", "l", false, "Monitor legacy migrations")
}

func initConfig() {
	// Configuration is handled by the config package
}

func runMigrationMonitor(cmd *cobra.Command, args []string) error {
	// Check for required organization flag
	if organization == "" {
		return fmt.Errorf("organization is required. Use --organization flag or set GHMM_GITHUB_ORGANIZATION environment variable")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override config with command line flags
	if organization != "" {
		cfg.GitHub.Organization = organization
	}
	if githubToken != "" {
		cfg.GitHub.Token = githubToken
	}
	if legacy {
		cfg.Migration.IsLegacy = legacy
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create GitHub client
	githubClient, err := api.NewGitHubClient(cfg.GitHub.Token, cfg.Migration.IsLegacy)
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Create migration service
	migrationService := services.NewMigrationService(githubClient)

	// Create UI dashboard
	dashboard := ui.NewDashboard()

	// Setup TUI application
	app := tview.NewApplication()
	grid := dashboard.SetupGrid()
	dashboard.SetupKeyboardNavigation(app, grid)

	// Start background data updates
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		// Initial load
		updateDashboard(ctx, migrationService, dashboard, cfg)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateDashboard(ctx, migrationService, dashboard, cfg)
			}
		}
	}()

	// Run the application
	return app.SetRoot(grid, true).SetFocus(grid).Run()
}

func updateDashboard(ctx context.Context, service services.MigrationService, dashboard *ui.Dashboard, cfg *config.Config) {
	summary, err := service.ListMigrations(ctx, cfg.GitHub.Organization, cfg.Migration.IsLegacy)
	if err != nil {
		// TODO: Add proper error handling/logging
		return
	}

	dashboard.UpdateData(summary)
}
