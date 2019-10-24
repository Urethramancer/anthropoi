package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasRemove options.
type CmdAliasRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdAliasRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
