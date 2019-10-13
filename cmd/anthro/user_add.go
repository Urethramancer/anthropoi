package main

import (
	"errors"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
)

// CmdUserAdd options.
type CmdUserAdd struct {
	opt.DefaultHelp
}

// Run add
func (cmd *CmdUserAdd) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db := anthropoi.New(host, port, username, password, name, ssl)
	m("%+v", db)
	err := db.Connect(name)
	if err != nil {
		return err
	}
	return nil
}
