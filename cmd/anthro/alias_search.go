package main

import (
	"encoding/json"
	"errors"

	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
)

// CmdAliasSearch options.
type CmdAliasSearch struct {
	opt.DefaultHelp
	Match string `placeholder:"KEYWORD" help:"Find aliases and targets containing keyword."`
	JSON  bool   `short:"j" long:"json" help:"Output in JSON format, suitable for import with the alias set subcommand."`
}

// Run get
func (cmd *CmdAliasSearch) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	list, err := db.SearchAliases(cmd.Match)
	if err != nil {
		return err
	}

	buf := stringer.New()
	if cmd.JSON {
		j, err := json.MarshalIndent(list, "", "\t")
		if err != nil {
			return err
		}
		buf.WriteString(string(j))
	} else {
		l := len(list.List)
		for i, a := range list.List {
			buf.WriteStrings("\t", a.Alias, " -> ", a.Target)
			if i < (l - 1) {
				buf.WriteString("\n")
			}
		}
	}

	m("%s", buf.String())
	return nil
}
