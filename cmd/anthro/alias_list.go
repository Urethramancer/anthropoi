package main

import (
	"encoding/json"
	"errors"

	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
)

// CmdAliasList options.
type CmdAliasList struct {
	opt.DefaultHelp
	Match string `placeholder:"KEYWORD" help:"Find aliases and targets containing keyword."`
	JSON  bool   `short:"j" long:"json" help:"Output in JSON format, suitable for import with the alias set subcommand."`
}

// Run get
func (cmd *CmdAliasList) Run(in []string) error {
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

	if cmd.JSON {
		j, err := json.MarshalIndent(list, "", "\t")
		if err != nil {
			return err
		}

		m("%s", string(j))
	} else {
		buf := stringer.New()
		l := len(list.List)
		for i, a := range list.List {
			buf.WriteStrings("\t", a.Alias, " -> ", a.Target)
			if i < (l - 1) {
				buf.WriteString("\n")
			}
		}
		m("%s", buf.String())
	}

	return nil
}
