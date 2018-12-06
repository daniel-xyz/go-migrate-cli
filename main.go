package migrate

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
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
		color.Red(err.Error())
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
			msgOnSuccess(nil, "Already up-to-date")
			return nil
		}

		msgOnSuccess(err, "Successfully migrated to latest version")
	case optionDown:
		err = m.Down()

		if err == migrate.ErrNoChange {
			msgOnSuccess(nil, "Already on lowest version")
			return nil
		}

		msgOnSuccess(err, "Successfully migrated to lowest version")
	case optionDrop:
		err = m.Drop()

		msgOnSuccess(err, "Successfully dropped tables and indexes")
	case optionForce:
		var v int

		fmt.Print("Migrate to which version? ")

		_, err = fmt.Fscanf(r, "%d", &v)

		if err != nil {
			return err
		}

		err = m.Force(v)

		msgOnSuccess(err, "Successfully migrated to forced version")
	case optionFullReset:
		err = m.Force(0)

		if err != nil {
			return err
		}

		err = m.Drop()

		msgOnSuccess(err, "Successfully dropped tables/indexes and forced version")
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
		color.Red(err.Error())

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
	version, _, err := m.Version()

	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	msg := fmt.Sprintf("Schema is at v%d", version)
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s %s\n\n", cyan("Status:"), msg)

	return nil
}

func msgOnSuccess(err error, msg string) {
	if err != nil {
		return
	}

	color.Green(msg)
}
