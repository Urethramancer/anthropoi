package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroupEdit options.
type CmdGroupEdit struct {
	opt.DefaultHelp
}

// Run edit
func (cmd *CmdGroupEdit) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

