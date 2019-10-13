package main

import (
	"errors"

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

	return nil
}
