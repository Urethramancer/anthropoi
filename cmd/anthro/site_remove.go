package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/signor/opt"
)

// CmdSiteRemove options.
type CmdSiteRemove struct {
	opt.DefaultHelp
	Site string `placeholder:"SITE" help:"Site or ID to remove. Users will need to be cleaned up with separate commands."`
}

// Run remove
func (cmd *CmdSiteRemove) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	id, err := strconv.ParseInt(cmd.Site, 10, 64)
	if err != nil {
		err = db.RemoveSiteByName(cmd.Site)
	} else {
		err = db.RemoveSite(id)
	}

	if err != nil {
		return err
	}

	m("Removed %s", cmd.Site)
	return nil
}
