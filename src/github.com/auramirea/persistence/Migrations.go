package persistence

import (
	"fmt"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

type DbMigration struct {
}

const DB_URL = "postgres://vagrant@192.168.10.11/defaultdb"
const MIGRATIONS_PATH = "db"

func (dbMigration *DbMigration) MigrationsUp() {
	// use synchronous versions of migration functions ...
	err, ok := migrate.UpSync(DB_URL, MIGRATIONS_PATH)
	if !ok {
		fmt.Println("Up sync error...", err, ok)
	}
}

func (dbMigration *DbMigration) MigrationsDown() {
	err, ok := migrate.DownSync(DB_URL, MIGRATIONS_PATH)
	if !ok {
		fmt.Println("Down sync error ...", err, ok)
	}
}
