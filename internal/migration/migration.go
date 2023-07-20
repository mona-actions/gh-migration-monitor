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

var migrations Migrations

func (m Migrations) FetchMigrations() {
	data := api.GetOrgMigrations()

	for _, migration := range data {
		switch migration.(map[string]interface{})["State"].(string) {
		case "QUEUED", "IN_PROGRESS", "SUCCEEDED":
			migrations.Queued = append(migrations.Queued, Migration{
				Id:             migration.(map[string]interface{})["Id"].(string),
				CreatedAt:      migration.(map[string]interface{})["CreatedAt"].(string),
				RepositoryName: migration.(map[string]interface{})["RepositoryName"].(string),
			})
		case "FAILED":
			migrations.Failed = append(migrations.Failed, Migration{
				Id:              migration.(map[string]interface{})["Id"].(string),
				CreatedAt:       migration.(map[string]interface{})["CreatedAt"].(string),
				FailureReason:   migration.(map[string]interface{})["FailureReason"].(string),
				RepositoryName:  migration.(map[string]interface{})["RepositoryName"].(string),
				MigrationLogUrl: migration.(map[string]interface{})["MigrationLogUrl"].(string),
			})
		}
	}
}
