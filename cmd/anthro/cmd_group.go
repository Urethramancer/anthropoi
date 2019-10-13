package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdGroup subcommands.
type CmdGroup struct {
	opt.DefaultHelp
	Add    CmdGroupAdd    `command:"add" help:"Add a new group."`
	Edit   CmdGroupEdit   `command:"edit" help:"Edit an existing group." aliases:"change"`
	List   CmdGroupList   `command:"list" help:"List groups." aliases:"ls"`
	Move   CmdGroupMove   `command:"move" help:"Move a group between sites." aliases:"mv"`
	Remove CmdGroupRemove `command:"remove" help:"Remove a group." aliases:"rm"`
}

// Run group
func (cmd *CmdGroup) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
