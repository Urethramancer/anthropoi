package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Urethramancer/anthropoi"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	db, err := openExisting()
	if err != nil {
		fmt.Printf("Error opening database: %s", err.Error)
		os.Exit(2)
	}

	defer db.Close()
	name := getenv("DB_NAME", "accounts")
	if !db.DatabaseExists(name) {
		fmt.Printf("No database. Setting up '%s' on '%s:%s'\n", getenv("DB_NAME", "accounts"), getenv("DB_HOST", "localhost"), getenv("DB_PORT", "5432"))
		db, err = openNew()
		defer db.Close()
		q := `CREATE DATABASE ` + name + `;`
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("Error creating database: %s\n", err.Error())
			os.Exit(2)
		}

		db, err = openExisting()
		if err != nil {
			fmt.Printf("Error opening database: %s", err.Error)
			os.Exit(2)
		}

		err = db.InitDatabase()
		if err != nil {
			fmt.Printf("Error initalising tables: %s", err.Error())
			os.Exit(2)
		}
	}
}

func getenv(key, alt string) string {
	s := os.Getenv(key)
	if s == "" {
		return alt
	}

	return s
}

func openExisting() (*anthropoi.DBM, error) {
	return anthropoi.OpenDB(
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		getenv("DB_NAME", "accounts"),
		getenv("DB_MODE", "disable"),
	)
}

func openNew() (*anthropoi.DBM, error) {
	return anthropoi.OpenDB(
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		"",
		getenv("DB_MODE", "disable"),
	)
}

func openDB(src string) *sql.DB {
	db, err := sql.Open("postgres", src)
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		os.Exit(2)
	}

	return db
}
