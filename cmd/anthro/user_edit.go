package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/signor/opt"
)

// CmdUserEdit options.
type CmdUserEdit struct {
	opt.DefaultHelp
	User string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
}

// Run edit
func (cmd *CmdUserEdit) Run(in []string) error {
	if cmd.Help || cmd.User == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect()
	if err != nil {
		return err
	}

	defer db.Close()
	id, err := strconv.ParseInt(cmd.User, 10, 64)
	if err != nil {
		m("String: %s", cmd.User)
	} else {
		m("Number: %d", id)
	}

	return nil
}
