package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/signor/opt"
)

// CmdUserRemove options.
type CmdUserRemove struct {
	opt.DefaultHelp
	User string `placeholder:"USER" help:"Username or ID to remove."`
}

// Run remove
func (cmd *CmdUserRemove) Run(in []string) error {
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
		err = db.DeleteUserByName(cmd.User)
	} else {
		err = db.DeleteUser(id)
	}

	if err != nil {
		return err
	}

	m("Removed user %s", cmd.User)
	return nil
}
