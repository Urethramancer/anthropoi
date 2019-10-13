package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdUser subcommands.
type CmdUser struct {
	Add    CmdUserAdd    `command:"add" help:"Add a new user."`
	Edit   CmdUserEdit   `command:"edit" help:"Edit an existing user." aliases:"ed,change"`
	List   CmdUserList   `command:"list" help:"List users." aliases:"ls"`
	Remove CmdUserRemove `command:"remove" help:"Remove a user." aliases:"rm"`
}

// Run user
func (cmd *CmdUser) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
