package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAlias subcommands.
type CmdAlias struct {
	List   CmdAliasList   `command:"list" help:"List aliases." aliases:"ls"`
	Set    CmdAliasSet    `command:"add" help:"Set a new or existing alias."`
	Search CmdAliasSearch `command:"search" help:"Search for an alias or target."`
	Remove CmdAliasRemove `command:"remove" help:"Remove an alias." aliases:"rm"`
}

// Run alias
func (cmd *CmdAlias) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
