package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdUser subcommands.
type CmdUser struct {
	List   CmdUserList   `command:"list" help:"List users." aliases:"ls,l"`
	Add    CmdUserAdd    `command:"add" help:"Add a new user. A password will be generated and displayed."`
	Edit   CmdUserEdit   `command:"edit" help:"Edit an existing user." aliases:"ed,change"`
	Reset  CmdUserReset  `command:"reset" help:"Reset password to a new random one."`
	Remove CmdUserRemove `command:"remove" help:"Remove a user." aliases:"rm"`
}

// Run user
func (cmd *CmdUser) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
