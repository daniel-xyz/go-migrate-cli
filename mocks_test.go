package migrate

import (
	"fmt"
	"github.com/golang-migrate/migrate"
)

type mockDB struct {
	migrationVersion int
}

type mockMigration struct {
	db                mockDB
	lastMigrationCall string
}

func (m *mockMigration) Up() error {
	m.lastMigrationCall = "Up"

	if m.db.migrationVersion == 1 {
		return migrate.ErrNoChange
	}

	return nil
}

func (m *mockMigration) Down() error {
	m.lastMigrationCall = "Down"

	if m.db.migrationVersion == 0 {
		return migrate.ErrNoChange
	}

	return nil
}

func (m *mockMigration) Drop() error {
	m.lastMigrationCall = "Drop"

	return nil
}

func (m *mockMigration) Force(v int) error {
	m.lastMigrationCall = fmt.Sprintf("Force(%d)", v)

	return nil
}

func (m *mockMigration) Version() (uint, bool, error) {
	m.lastMigrationCall = "Version"

	return uint(m.db.migrationVersion), false, nil
}

func (m *mockMigration) Close() (error, error) {
	m.lastMigrationCall = "Close"

	return nil, nil
}
