package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroupRemove options.
type CmdGroupRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdGroupRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

