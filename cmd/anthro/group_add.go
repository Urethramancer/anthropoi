package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroupAdd options.
type CmdGroupAdd struct {
	opt.DefaultHelp
}

// Run add
func (cmd *CmdGroupAdd) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

