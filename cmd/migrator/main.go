package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var host, port, name, user, password, migrationPath, migrationTable string
	flag.StringVar(&host, "host", "127.0.0.1", "Host of the database")
	flag.StringVar(&port, "port", "5432", "Port of the database")
	flag.StringVar(&name, "name", "postgres", "Name of the database")
	flag.StringVar(&user, "user", "postgres", "User of the database")
	flag.StringVar(&password, "password", "postgres", "Password of the database")
	flag.StringVar(&migrationPath, "migration", "", "Path to the migration files")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "Table name to store the migration history")
	flag.Parse()

	if host == "" || port == "" || name == "" || user == "" || password == "" || migrationPath == "" || migrationTable == "" {
		panic("host, port, name, user, password, migration-path and migration-table are required")
	}

	driverURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, name)
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
