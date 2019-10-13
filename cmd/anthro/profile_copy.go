package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfileCopy options.
type CmdProfileCopy struct {
	opt.DefaultHelp
}

// Run copy
func (cmd *CmdProfileCopy) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

