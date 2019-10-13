package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfileAdd options.
type CmdProfileAdd struct {
	opt.DefaultHelp
	User string `placeholder:"USER" help:"Username or ID to add a profile to."`
}

// Run add
func (cmd *CmdProfileAdd) Run(in []string) error {
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
