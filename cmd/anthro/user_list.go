package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdUserList options.
type CmdUserList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdUserList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
