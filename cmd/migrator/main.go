package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storage", "", "Path to the storage file")
	flag.StringVar(&migrationPath, "migration", "", "Path to the migration files")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "Table name to store the migration history")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}
	if migrationPath == "" {
		panic("migration path is required")
	}

	driverURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", "postgres", "127.0.0.1", "5432", "postgres")
	m, err := migrate.New(
		"file://"+migrationPath,
		driverURL,
		//TODO: set the table name for the migration history
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("Migrations applied")
}
