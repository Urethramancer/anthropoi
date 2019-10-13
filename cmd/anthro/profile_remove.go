package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfileRemove options.
type CmdProfileRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdProfileRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

