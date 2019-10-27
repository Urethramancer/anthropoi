package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAlias subcommands.
type CmdAlias struct {
	List   CmdAliasList   `command:"list" help:"List or search for aliases and targets." aliases:"ls,l"`
	Set    CmdAliasSet    `command:"add" help:"Set a new or existing alias." aliases:"a"`
	Import CmdAliasImport `command:"import" help:"Import aliases from a JSON file." aliases:"imp,i"`
	Remove CmdAliasRemove `command:"remove" help:"Remove an alias." aliases:"rm"`
}

// Run alias
func (cmd *CmdAlias) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
