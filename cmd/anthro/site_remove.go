package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdSiteRemove options.
type CmdSiteRemove struct {
	opt.DefaultHelp
}

// Run remove
func (cmd *CmdSiteRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	return nil
}
