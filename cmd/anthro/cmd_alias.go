package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAlias subcommands.
type CmdAlias struct {
	Set    CmdAliasSet    `command:"add" help:"Set a new or existing alias."`
	Search CmdAliasSearch `command:"search" help:"Search for an alias or target." aliases:"list,ls"`
	Remove CmdAliasRemove `command:"remove" help:"Remove an alias." aliases:"rm"`
}

// Run alias
func (cmd *CmdAlias) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
