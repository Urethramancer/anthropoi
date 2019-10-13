package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// Cmdinit options.
type CmdInit struct {
	opt.DefaultHelp
	Drop bool `short:"D" long:"drop" help:"Drop existing database or tables."`
}

// Run init
func (cmd *CmdInit) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect()
	if err != nil {
		e("Error opening database: %s", err.Error())
		return err
	}

	if cmd.Drop {
		err := db.Drop(name)
		if err != nil {
			e("Couldn't drop database: %s", err.Error())
			return err
		}
	}

	defer db.Close()
	if !db.DatabaseExists(name) {
		m("No database. Setting up '%s' on '%s:%s'",
			getenv("DB_NAME", name),
			getenv("DB_HOST", host),
			getenv("DB_PORT", port),
		)

		err = db.Connect("")
		if err != nil {
			e("Error opening database: %s", err.Error())
			return err
		}

		defer db.Close()
		err = db.Create(name)
		if err != nil {
			e("Error creating database: %s", err.Error())
			return err
		}

	}

	err = db.Connect(name)
	if err != nil {
		e("Error opening database: %s", err.Error())
		return err
	}

	err = db.InitDatabase()
	if err != nil {
		e("Error initalising database: %s", err.Error())
		return err
	}

	return nil
}
