package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasRemove options.
type CmdAliasRemove struct {
	opt.DefaultHelp
	Alias string `placeholder:"ALIAS" help:"Alias to remove."`
}

// Run remove
func (cmd *CmdAliasRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	err = db.RemoveAlias(cmd.Alias)
	if err != nil {
		return err
	}

	m("%s removed.", cmd.Alias)
	return nil
}
