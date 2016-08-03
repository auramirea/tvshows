package persistence

import (
	"github.com/mattes/migrate/migrate"
	"fmt"
)
type DbMigration struct {
}

const DB_URL = "postgres://vagrant@192.168.10.11/defaultdb"
const MIGRATIONS_PATH = "db"

func (dbMigration *DbMigration) MigrationsUp() {
	// use synchronous versions of migration functions ...
	err, ok := migrate.UpSync(DB_URL, MIGRATIONS_PATH)
	if !ok {
		fmt.Println("Oh no ...")
		// do sth with allErrors slice
		for e := range err {
			fmt.Println(e)
		}
	}
}

func (dbMigration *DbMigration) MigrationsDown() {
	err, ok := migrate.DownSync(DB_URL, MIGRATIONS_PATH)
	if !ok {
		fmt.Println("Down sync error ...")
		// do sth with allErrors slice
		for e := range err {
			fmt.Println(e)
		}
	}
}

