package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdUserEdit options.
type CmdUserEdit struct {
	opt.DefaultHelp
}

// Run edit
func (cmd *CmdUserEdit) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
