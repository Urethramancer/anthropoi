package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfileList options.
type CmdProfileList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdProfileList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

