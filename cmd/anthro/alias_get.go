package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasGet options.
type CmdAliasGet struct {
	opt.DefaultHelp
}

// Run get
func (cmd *CmdAliasGet) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
