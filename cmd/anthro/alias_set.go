package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdAliasSet options.
type CmdAliasSet struct {
	opt.DefaultHelp
	Alias  string `placeholder:"ALIAS" help:"Alias to create. Must be @ a domain in the database."`
	Target string `placeholder:"TARGET" help:"Address to forward to. This will be resolved to an actual account if you specify an alias."`
	JSON   string `short:"j" long:"json" help:"Specify JSON file containing alias:target dictionary pairs of aliases to add."`
}

// Run set
func (cmd *CmdAliasSet) Run(in []string) error {
	if cmd.Help || cmd.Target == "" {
		if cmd.JSON == "" {
			return errors.New(opt.ErrorUsage)
		}
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	t, err := db.GetAlias(cmd.Target)
	if err != nil {
		return err
	}

	err = db.SetAlias(cmd.Alias, cmd.Target)
	if err != nil {
		return err
	}

	m("%s -> %s", cmd.Alias, t)
	return nil
}
