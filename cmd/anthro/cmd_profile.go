package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// CmdProfile subcommands.
type CmdProfile struct {
	opt.DefaultHelp
	List      CmdProfileList      `command:"list" help:"List profiles." aliases:"ls"`
	Add       CmdProfileAdd       `command:"add" help:"Add new profile to a user."`
	SetGroups CmdProfileSetGroups `command:"setgroups" help:"Set groups for a profile." aliases:"sg"`
	Remove    CmdProfileRemove    `command:"remove" help:"Remove a profile." aliases:"rm"`
	Copy      CmdProfileCopy      `command:"copy" help:"Copy a profile from one user to another." aliases:"cp"`
}

// Run profile
func (cmd *CmdProfile) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
