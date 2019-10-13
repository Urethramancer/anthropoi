package main

import (
	"errors"
	"strconv"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
)

// CmdUserReset options.
type CmdUserReset struct {
	opt.DefaultHelp
	User string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
	Len  int    `short:"l" long:"length" help:"Length of password." default:"14"`
	Cost int    `short:"c" long:"cost" help:"Cost of hashing algorithm. Tweak this to around 100ms per hash." default:"11"`
}

// Run reset
func (cmd *CmdUserReset) Run(in []string) error {
	if cmd.Help || cmd.User == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect()
	if err != nil {
		return err
	}

	defer db.Close()
	id, err := strconv.ParseInt(cmd.User, 10, 64)
	var u *anthropoi.User
	if err != nil {
		u, err = db.GetUserByName(cmd.User)
	} else {
		u, err = db.GetUser(id)
	}
	if err != nil {
		return err
	}

	pw := anthropoi.GenString(cmd.Len)
	u.SetPassword(pw, cmd.Cost)
	m("Changed password for %s%s(%d)%s to %s%s%s", ansi.Blue, u.Usermame, u.ID, ansi.Normal, ansi.Green, pw, ansi.Normal)
	return nil
}
