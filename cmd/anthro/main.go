package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/anthropoi"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	db := anthropoi.New(
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		"",
		getenv("DB_MODE", "disable"),
	)

	err = db.Connect("")
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		os.Exit(2)
	}

	defer db.Close()
	name := getenv("DB_NAME", "accounts")
	if !db.DatabaseExists(name) {
		fmt.Printf("No database. Setting up '%s' on '%s:%s'\n",
			getenv("DB_NAME", "accounts"),
			getenv("DB_HOST", "localhost"),
			getenv("DB_PORT", "5432"),
		)

		err = db.Connect("")
		if err != nil {
			fmt.Printf("Error opening database: %s\n", err.Error())
			os.Exit(2)
		}

		println("Opened " + db.ConnectionString())
		defer db.Close()
		err = db.Create(getenv("DB_NAME", "accounts"))
		if err != nil {
			fmt.Printf("Error creating database: %s\n", err.Error())
			os.Exit(2)
		}

		err = db.Connect(name)
		if err != nil {
			fmt.Printf("Error opening database: %s\n", err.Error())
			os.Exit(2)
		}

		err = db.InitDatabase(name)
		if err != nil {
			fmt.Printf("Error initalising database: %s\n", err.Error())
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
