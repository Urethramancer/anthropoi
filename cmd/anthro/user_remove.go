package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdUserRemove options.
type CmdUserRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdUserRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
