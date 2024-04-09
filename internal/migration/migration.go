package migration

import (
	"github.com/mona-actions/gh-migration-monitor/internal/api"
)

type Migration struct {
	Id              string
	CreatedAt       string
	FailureReason   string
	RepositoryName  string
	State           string
	MigrationLogUrl string
}

type Migrations struct {
	Queued      []Migration
	In_Progress []Migration
	Succeeded   []Migration
	Failed      []Migration
	Log         []Migration
}

func GetMigrations() Migrations {
	var migrations Migrations

	// Get the data from the API
	data := api.GetOrgMigrations()

	for _, migration := range data {
		repo := Migration{
			Id:             migration["Id"],
			CreatedAt:      migration["CreatedAt"],
			RepositoryName: migration["RepositoryName"],
		}

		switch migration["State"] {
		case "QUEUED", "WAITING":
			migrations.Queued = append(migrations.Queued, repo)
		case "IN_PROGRESS", "PREPARING", "PENDING", "MAPPING", "ARCHIVE_UPLOADED", "CONFLICTS", "READY", "IMPORTING":
			migrations.In_Progress = append(migrations.In_Progress, repo)
		case "SUCCEEDED", "UNLOCKED", "IMPORTED":
			migrations.Succeeded = append(migrations.Succeeded, repo)
		case "FAILED", "FAILED_IMPORT":
			repo.FailureReason = migration["FailureReason"]
			repo.MigrationLogUrl = migration["MigrationLogUrl"]

			migrations.Failed = append(migrations.Failed, repo)
		}
	}

	return migrations
}
