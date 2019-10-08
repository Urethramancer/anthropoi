package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Urethramancer/anthropoi"
	_ "github.com/lib/pq"
)

func main() {
	db := openExisting()
	defer db.Close()
	name := getenv("DB_NAME", "accounts")
	if !anthropoi.DatabaseExists(db, name) {
		fmt.Printf("No database. Setting up '%s' on '%s:%s'\n", getenv("DB_NAME", "accounts"), getenv("DB_HOST", "localhost"), getenv("DB_PORT", "5432"))
		db = openNew()
		defer db.Close()
		q := `CREATE DATABASE ` + name + `;`
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("Error creating database: %s\n", err.Error())
			os.Exit(2)
		}

		db = openExisting()
		err = anthropoi.InitDatabase(db)
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

func openExisting() *sql.DB {
	src := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		getenv("DB_NAME", "accounts"),
		getenv("DB_MODE", "disable"),
	)
	return openDB(src)
}

func openNew() *sql.DB {
	src := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		getenv("DB_MODE", "disable"),
	)
	return openDB(src)
}

func openDB(src string) *sql.DB {
	db, err := sql.Open("postgres", src)
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		os.Exit(2)
	}

	return db
}
