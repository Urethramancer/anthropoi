package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroupMove options.
type CmdGroupMove struct {
	opt.DefaultHelp
}

// Run move
func (cmd *CmdGroupMove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}

