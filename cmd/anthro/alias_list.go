package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasList options.
type CmdAliasList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdAliasList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	if !db.GetFlag("mailmode") {
		m("Database is not set up for mail mode.")
		return nil
	}

	return nil
}
