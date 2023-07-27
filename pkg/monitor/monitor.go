package monitor

import (
	"fmt"

	"github.com/mona-actions/gh-migration-monitor/internal/migration"
)

func Organization() {
	var migrations migration.Migrations

	migrations.FetchMigrations()

	fmt.Println("Queued Migrations:" + fmt.Sprint(migrations.Queued))
	fmt.Println("In Progress Migrations:" + fmt.Sprint(migrations.In_Progress))
	fmt.Println("Succeeded Migrations:" + fmt.Sprint(migrations.Succeeded))
}
