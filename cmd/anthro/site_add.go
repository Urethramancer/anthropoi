package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdSiteAdd options.
type CmdSiteAdd struct {
	opt.DefaultHelp
	Site string `placeholder:"SITE" help:"Site or domain name to add."`
}

// Run add
func (cmd *CmdSiteAdd) Run(in []string) error {
	if cmd.Help || cmd.Site == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	id, err := db.AddSite(cmd.Site)
	if err != nil {
		return err
	}

	m("%s added with ID %d.", cmd.Site, id)
	return nil
}
