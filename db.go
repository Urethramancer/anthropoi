package anthropoi

import (
	"database/sql"
	"fmt"
)

// DBM is a DB manager for user accounts and groups.
type DBM struct {
	*sql.DB
}

func OpenDB(host, port, user, password, name, mode string) (*DBM, error) {
	dbm := DBM{}
	if mode != "enable" && mode != "disable" {
		mode = "disable"
	}
	src := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, name, mode,
	)

	var err error
	dbm.DB, err = sql.Open("postgres", src)
	if err != nil {
		return nil, err
	}

	return &dbm, nil
}

// InitDatabase creates the tables, functions and triggers required for the full account system.
func (db *DBM) InitDatabase() error {
	var err error
	_, err = db.Exec(functionDefinitions)
	if err != nil {
		return err
	}

	_, err = db.Exec(userTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(profileTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(groupTables)
	if err != nil {
		return err
	}

	return nil
}

// DatabaseExists checks for the existence of the actual database.
func (db *DBM) DatabaseExists(name string) bool {
	q := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '` + name + `');`
	row := db.QueryRow(q)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}
