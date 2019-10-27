package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
)

// CmdAliasImport options.
type CmdAliasImport struct {
	opt.DefaultHelp
	File string `placeholder:"FILE" help:"JSON file containing alias-target pairs."`
}

// Run import
func (cmd *CmdAliasImport) Run(in []string) error {
	if cmd.Help || cmd.File == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	data, err := ioutil.ReadFile(cmd.File)
	if err != nil {
		return err
	}

	var aliases anthropoi.Aliases
	err = json.Unmarshal(data, &aliases)
	if err != nil {
		return err
	}

	for _, a := range aliases.List {
		err = db.SetAlias(a.Alias, a.Target)
		if err != nil {
			m("Failed to set %s", a.Alias)
		} else {
			m("%s -> %s", a.Alias, a.Target)
		}
	}
	return nil
}
