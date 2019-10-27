package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdSite subcommands.
type CmdSite struct {
	opt.DefaultHelp
	List   CmdSiteList   `command:"list" help:"List or search for sites." aliases:"ls,l"`
	Add    CmdSiteAdd    `command:"add" help:"Add new site."`
	Remove CmdSiteRemove `command:"remove" help:"Remove a site." aliases:"rm"`
}

// Run sites
func (cmd *CmdSite) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
