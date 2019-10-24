package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasSet options.
type CmdAliasSet struct {
	opt.DefaultHelp
	Alias  string `placeholder:"" help:""`
	Target string `placeholder:"" help:""`
}

// Run set
func (cmd *CmdAliasSet) Run(in []string) error {
	if cmd.Help || cmd.Target == "" {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
