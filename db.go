package anthropoi

import (
	"database/sql"
	"strings"

	"github.com/Urethramancer/signor/stringer"
)

// DBM is a DB manager for user accounts and groups.
type DBM struct {
	*sql.DB
	host     string
	port     string
	user     string
	password string
	name     string
	mode     string
	cs       *stringer.Stringer
}

// New DBM setup.
func New(host, port, user, password, mode string) *DBM {
	dbm := &DBM{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		mode:     mode,
		cs:       stringer.New(),
	}

	if host == "" {
		dbm.host = "localhost"
	}

	if port == "" {
		dbm.port = "5432"
	}

	if user == "" {
		dbm.user = "postgres"
	}

	if mode == "" {
		dbm.mode = "disable"
	}
	return dbm
}

// OpenDB and set the pointer in the DBM struct.
func (db *DBM) Connect(name string) error {
	var err error
	db.name = name
	db.DB, err = sql.Open("postgres", db.ConnectionString())
	if err != nil {
		return err
	}

	return nil
}

func (db *DBM) ConnectionString() string {
	db.cs.Reset()
	db.cs.WriteStrings(
		"host=", db.host, " ",
		"port=", db.port, " ",
		"user=", db.user, " ",
	)

	if db.password != "" {
		db.cs.WriteStrings("password=", db.password, " ")
	}

	if db.name != "" {
		db.cs.WriteStrings("dbname=", db.name, " ")
	}

	if db.mode == "enable" {
		db.cs.WriteString("sslmode=enable")
	} else {
		db.cs.WriteString("sslmode=disable")
	}

	return db.cs.String()
}

// Create the database and retain the name.
func (db *DBM) Create(name string) error {
	if name == "" {
		name = DefaultName
	}
	q := strings.Replace(databaseDefinitions, "{NAME}", name, 1)
	_, err := db.Exec(q)
	if err == nil {
		db.name = name
	}
	return err
}

// InitDatabase creates the tables, functions and triggers required for the full account system.
func (db *DBM) InitDatabase() error {
	var err error
	if db.DB != nil {
		err = db.Close()
		if err != nil {
			return err
		}
	}

	db.DB, err = sql.Open("postgres", db.ConnectionString())
	if err != nil {
		return err
	}

	_, err = db.Exec(databaseTriggers)
	if err != nil {
		return err
	}

	_, err = db.Exec(userTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(groupTables)
	return err
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

// Drop a named database.
func (db *DBM) Drop(name string) error {
	var err error
	q := `SELECT pid, pg_terminate_backend(pid) 
	FROM pg_stat_activity 
	WHERE datname = '` + name + `' AND pid <> pg_backend_pid();`
	_, err = db.Exec(q)
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + name + ";")
	return err
}
