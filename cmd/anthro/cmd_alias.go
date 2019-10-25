package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAlias subcommands.
type CmdAlias struct {
	Set    CmdAliasSet    `command:"add" help:"Set a new or existing alias."`
	List   CmdAliasList   `command:"list" help:"List or search for aliases and targets." aliases:"ls"`
	Remove CmdAliasRemove `command:"remove" help:"Remove an alias." aliases:"rm"`
}

// Run alias
func (cmd *CmdAlias) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
