package migration

import "github.com/mona-actions/gh-migration-monitor/internal/api"

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

func (m *Migrations) FetchMigrations() {
	data := api.GetOrgMigrations()

	for _, migration := range data {
		repo := Migration{
			Id:             migration["Id"],
			CreatedAt:      migration["CreatedAt"],
			RepositoryName: migration["RepositoryName"],
		}

		switch migration["State"] {
		case "QUEUED":
			m.Queued = append(m.Queued, repo)
		case "IN_PROGRESS":
			m.In_Progress = append(m.In_Progress, repo)
		case "SUCCEEDED":
			m.Succeeded = append(m.Succeeded, repo)
		case "FAILED":
			repo.FailureReason = migration["FailureReason"]
			repo.MigrationLogUrl = migration["MigrationLogUrl"]

			m.Failed = append(m.Failed, repo)
		}
	}
}
