package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfileSetGroups options.
type CmdProfileSetGroups struct {
	opt.DefaultHelp
}

// Run edit
func (cmd *CmdProfileSetGroups) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
