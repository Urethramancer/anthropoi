package main

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
	"golang.org/x/crypto/ssh/terminal"
)

// CmdUserReset options.
type CmdUserReset struct {
	opt.DefaultHelp
	User    string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
	Len     int    `short:"l" long:"length" help:"Length of password." default:"14"`
	Length  int    `short:"L" long:"length" help:"Length of password. Minimum allowed is 12 characters." default:"14"`
	Cost    int    `short:"c" long:"cost" help:"Cost of bcrypt hashing algorithm. Tweak this to around 100ms per hash." default:"11"`
	Rounds  int    `short:"r" long:"rounds" help:"Number of rounds to hash SHA512-CRYPT." default:"50000"`
	Dovecot bool   `short:"d" long:"dovecot" help:"Generate a Dovecot-compatible password using SHA512-CRYPT, rather than the default bcrypt hash."`
	Ask     bool   `short:"a" long:"ask" help:"Ask for a password to set instead of generating one. This is the most secure option."`
}

// Run reset
func (cmd *CmdUserReset) Run(in []string) error {
	if cmd.Help || cmd.User == "" {
		return errors.New(opt.ErrorUsage)
	}

	// Enforce some sane minimums
	if cmd.Length < 12 {
		cmd.Length = 12
	}

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

	var pw string
	if cmd.Ask {
		fmt.Printf("Password: ")
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		println()
		if len(pass) == 0 {
			return errors.New("no password entered")
		}

		pw = string(pass)
	} else {
		pw = anthropoi.GenString(cmd.Length)
	}
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

	if !cmd.Ask {
		m("Changed password for %s%s(%d)%s to %s%s%s", ansi.Blue, u.Username, u.ID, ansi.Normal, ansi.Green, pw, ansi.Normal)
	}
	return nil
}
