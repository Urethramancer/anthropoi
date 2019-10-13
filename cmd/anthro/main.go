package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	_ "github.com/lib/pq"
)

// Options holds all the tool commands.
var Options struct {
	opt.DefaultHelp
	User  CmdUser  `command:"user" help:"User management."`
	Group CmdGroup `command:"group" help:"Group management."`
}

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

	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	err = a.RunCommand(false)
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}
}

func getenv(key, alt string) string {
	s := os.Getenv(key)
	if s == "" {
		return alt
	}

	return s
}
