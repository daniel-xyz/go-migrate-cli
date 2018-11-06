package migrate

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // needed for golang-migrate
	_ "github.com/lib/pq"                             // needed for golang-migrate
	"github.com/manifoldco/promptui"
	"os"
)

var m *migrate.Migrate

const (
	optionUp      = "Up"
	optionDown    = "Down"
	optionDrop    = "Drop all tables & indexes"
	optionForce   = "Force specific version"
	optionNothing = "Do nothing"
)

type settings struct {
	connection       *sql.DB
	dbName           string
	migrationsFolder string
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// CLI starts the cli with the specified settings
func CLI(connection *sql.DB, dbName string, migrationsFolder string) {
	s := settings{connection, dbName, migrationsFolder}
	m = getMigrateInstance(s)
	defer m.Close()

	printVersion()

	result := promptSelect()

	switch result {
	case optionUp:
		err := m.Up()

		if err == migrate.ErrNoChange {
			fmt.Println("Already up-to-date.")
			break
		}

		checkErr(err)
		fmt.Println("\nMigration successful.")
	case optionDown:
		err := m.Down()

		if err == migrate.ErrNoChange {
			break
		}

		checkErr(err)
		fmt.Println("\nMigration successful.")
	case optionDrop:
		err := m.Drop()

		checkErr(err)
		fmt.Println("\nMigration successful.")
	case optionForce:
		var v int
		fmt.Print("Version? ")

		_, err := fmt.Scanf("%d", &v)
		checkErr(err)

		err = m.Force(v)
		checkErr(err)

		fmt.Println("\nMigration successful.")
	case optionNothing:
		os.Exit(0)
	}

	printVersion()
}

func getMigrateInstance(s settings) *migrate.Migrate {
	driver, err := postgres.WithInstance(s.connection, &postgres.Config{})
	checkErr(err)

	m, err = migrate.NewWithDatabaseInstance("file://"+s.migrationsFolder, s.dbName, driver)
	checkErr(err)

	return m
}

func promptSelect() string {
	prompt := promptui.Select{
		Label: "Choose option",
		Items: []string{optionUp, optionDown, optionDrop, optionForce, optionNothing},
	}

	_, result, err := prompt.Run()
	checkErr(err)

	return result
}

func printVersion() {
	v, _, err := m.Version()

	if err == migrate.ErrNilVersion {
		fmt.Println("\nNo migrations have been done yet. ")
	}

	if err != migrate.ErrNilVersion {
		checkErr(err)
	}

	fmt.Printf("Schema is at v%d.\n\n", v)
}
