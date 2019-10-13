package main

import (
	"errors"
	"os"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
)

// Cmdinit options.
type CmdInit struct {
	opt.DefaultHelp
	Drop     bool   `short:"D" long:"drop" help:"Drop existing database or tables."`
	Host     string `short:"H" long:"host" help:"Host to connect to." default:"localhost"`
	Port     string `short:"p" long:"port" help:"Port to connect to." default:"5432"`
	User     string `short:"u" long:"user" help:"User to connect as." default:"postgres"`
	Password string `short:"P" long:"password" help:"Password for that user. Nay be left out if PostgreSQL is configured for other authentication methods."`
	Name     string `short:"n" long:"name" help:"Name of the database to create." default:"accounts"`
	SSL      bool   `short:"s" long:"ssl" help:"Require SSL to connect."`
}

// Run init
func (cmd *CmdInit) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	ssl := "disable"
	if cmd.SSL {
		ssl = "enable"
	}

	var err error
	db := anthropoi.New(
		getenv("DB_HOST", cmd.Host),
		getenv("DB_PORT", cmd.Port),
		getenv("DB_USER", cmd.User),
		getenv("DB_PASSWORD", cmd.Password),
		"",
		getenv("DB_MODE", ssl),
	)

	err = db.Connect("")
	if err != nil {
		e("Error opening database: %s", err.Error())
		os.Exit(2)
	}

	if cmd.Drop {
		q := "DROP DATABASE IF EXISTS " + getenv("DB_NAME", cmd.Name) + ";"
		_, err := db.Exec(q)
		if err != nil {
			e("Couldn't drop database: %s", err.Error())
			os.Exit(2)
		}
	}

	defer db.Close()
	name := getenv("DB_NAME", anthropoi.DefaultName)
	if !db.DatabaseExists(name) {
		m("No database. Setting up '%s' on '%s:%s'",
			getenv("DB_NAME", name),
			getenv("DB_HOST", cmd.Host),
			getenv("DB_PORT", cmd.Port),
		)

		err = db.Connect("")
		if err != nil {
			e("Error opening database: %s", err.Error())
			os.Exit(2)
		}

		defer db.Close()
		err = db.Create(name)
		if err != nil {
			e("Error creating database: %s", err.Error())
			os.Exit(2)
		}

	}

	err = db.Connect(name)
	if err != nil {
		e("Error opening database: %s", err.Error())
		os.Exit(2)
	}

	err = db.InitDatabase()
	if err != nil {
		e("Error initalising database: %s", err.Error())
		os.Exit(2)
	}

	return nil
}
