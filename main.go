package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // needed for golang-migrate
	_ "github.com/lib/pq"                             // needed for golang-migrate
	"github.com/manifoldco/promptui"
	"io"
	"os"
)

const (
	optionUp        = "Up - all versions"
	optionDown      = "Down - all versions"
	optionDrop      = "Drop - all tables/indexes"
	optionForce     = "Force - specific version"
	optionFullReset = "Reset - force first version & drop all tables/indexes"
	optionNothing   = "Do nothing - exit"
)

type migrationInstance interface {
	Up() error
	Down() error
	Drop() error
	Force(int) error
	Version() (uint, bool, error)
	Close() (error, error)
}

type settings struct {
	connection       *sql.DB
	dbName           string
	migrationsFolder string
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

// CLI starts the cli with the specified settings
func CLI(connection *sql.DB, dbName string, migrationsFolder string) {
	s := settings{connection, dbName, migrationsFolder}
	m, err := getMigrateInstance(s)
	defer m.Close()
	handleError(err)

	err = printVersion(m)
	handleError(err)

	startPrompt(m)

	err = printVersion(m)
	handleError(err)

	os.Exit(0)
}

func executeOption(r io.Reader, m migrationInstance, optionKey string) error {
	var err error

	switch optionKey {
	case optionUp:
		err = m.Up()

		if err == migrate.ErrNoChange {
			return errors.New("already up-to-date")
		}
	case optionDown:
		err = m.Down()

		if err == migrate.ErrNoChange {
			return errors.New("already on lowest possible version")
		}
	case optionDrop:
		err = m.Drop()
	case optionForce:
		var v int

		fmt.Print("Migrate to which version? ")

		_, err = fmt.Fscanf(r, "%d", &v)

		if err != nil {
			return err
		}

		err = m.Force(v)
	case optionFullReset:
		err = m.Force(0)

		if err != nil {
			return err
		}

		err = m.Drop()
	case optionNothing:
		return nil
	}

	return err
}

func getMigrateInstance(s settings) (migrationInstance, error) {
	driver, err := postgres.WithInstance(s.connection, &postgres.Config{})

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+s.migrationsFolder, s.dbName, driver)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func startPrompt(m migrationInstance) {
	selectedOption, err := promptSelect()

	if err = executeOption(os.Stdin, m, selectedOption); err != nil {
		fmt.Println(err.Error())

		startPrompt(m)
	}
}

func promptSelect() (string, error) {
	prompt := promptui.Select{
		Label: "Choose option",
		Items: []string{optionUp, optionDown, optionDrop, optionForce, optionFullReset, optionNothing},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func printVersion(m migrationInstance) error {
	v, _, err := m.Version()

	if err == migrate.ErrNilVersion {
		fmt.Println("\nNo migrations have been done yet. ")
	}

	if err != migrate.ErrNilVersion {
		return err
	}

	fmt.Printf("Schema is at v%d.\n\n", v)

	return nil
}
