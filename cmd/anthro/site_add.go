package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdSiteAdd options.
type CmdSiteAdd struct {
	opt.DefaultHelp
}

// Run add
func (cmd *CmdSiteAdd) Run(in []string) error {
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
