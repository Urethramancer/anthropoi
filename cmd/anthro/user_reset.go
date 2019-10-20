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
	User    string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
	Len     int    `short:"l" long:"length" help:"Length of password." default:"14"`
	Cost    int    `short:"c" long:"cost" help:"Cost of bcrypt hashing algorithm. Tweak this to around 100ms per hash." default:"11"`
	Dovecot bool   `short:"d" long:"dovecot" help:"Generate a Dovecot-compatible password using SHA512-CRYPT, rather than the default bcrypt hash."`
	Rounds  int    `short:"r" long:"rounds" help:"Number of rounds to hash SHA512-CRYPT." default:"50000"`
}

// Run reset
func (cmd *CmdUserReset) Run(in []string) error {
	if cmd.Help || cmd.User == "" {
		return errors.New(opt.ErrorUsage)
	}

	// Enforce some sane minimums
	if cmd.Cost < 10 {
		cmd.Cost = 10
	}

	if cmd.Rounds < 10000 {
		cmd.Rounds = 10000
	}

	db, err := connect(name)
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
	if cmd.Dovecot {
		u.SetDovecotPassword(pw, cmd.Rounds)
	} else {
		err = u.SetPassword(pw, cmd.Cost)
		if err != nil {
			return err
		}
	}

	err = db.SaveUser(u)
	if err != nil {
		return err
	}

	m("Changed password for %s%s(%d)%s to %s%s%s", ansi.Blue, u.Usermame, u.ID, ansi.Normal, ansi.Green, pw, ansi.Normal)
	return nil
}
