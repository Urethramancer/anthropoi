package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroupList options.
type CmdGroupList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdGroupList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

