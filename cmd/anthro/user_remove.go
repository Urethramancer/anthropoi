package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/anthropoi"
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

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	var u *anthropoi.User
	id, err := strconv.ParseInt(cmd.User, 10, 64)
	if err != nil {
		u, err = db.GetUserByName(cmd.User)
		if err != nil {
			return err
		}

		err = db.RemoveUserByName(cmd.User)
	} else {
		u, err = db.GetUser(id)
		if err != nil {
			return err
		}

		err = db.RemoveUser(id)
	}

	if err != nil {
		return err
	}

	if db.GetFlag("mailmode") {
		err = db.RemoveAliases(u.Username)
		if err != nil {
			return err
		}
	}

	m("Removed user %s", cmd.User)
	return nil
}
