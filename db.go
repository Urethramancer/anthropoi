package anthropoi

import (
	"database/sql"

	"github.com/Urethramancer/anthropoi/groups"
	"github.com/Urethramancer/anthropoi/profiles"
	"github.com/Urethramancer/anthropoi/users"
)

// InitDatabase creates the tables, functions and triggers required for the full account system.
func InitDatabase(db *sql.DB) error {
	_, err := db.Exec(functionDefinitions)
	if err != nil {
		return err
	}

	err = users.InitTables(db)
	if err != nil {
		return err
	}

	err = profiles.InitTables(db)
	if err != nil {
		return err
	}

	err = groups.InitTables(db)
	if err != nil {
		return err
	}

	return nil
}

// DatabaseExists checks for the existence of the actual database.
func DatabaseExists(db *sql.DB, name string) bool {
	q := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '` + name + `');`
	row := db.QueryRow(q)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func TableExists(name string) bool {
	return false
}
