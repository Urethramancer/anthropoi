package main

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
	"golang.org/x/crypto/ssh/terminal"
)

// CmdUserAdd options.
type CmdUserAdd struct {
	opt.DefaultHelp
	Name    string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
	Email   string `short:"e" long:"email" help:"Optional e-mail"`
	First   string `short:"f" long:"firstname" help:"Optional first name (the one displayed first - may be family name for some regions)."`
	Last    string `short:"l" long:"lastname" help:"Optional last name."`
	Length  int    `short:"L" long:"length" help:"Length of password. Minimum allowed is 12 characters." default:"14"`
	Cost    int    `short:"c" long:"cost" help:"Cost of hashing algorithm. Tweak this to around 100ms per hash." default:"11"`
	Rounds  int    `short:"r" long:"rounds" help:"Number of rounds to hash SHA512-CRYPT." default:"50000"`
	Dovecot bool   `short:"d" long:"dovecot" help:"Generate a Dovecot-compatible password using SHA512-CRYPT, rather than the default bcrypt hash."`
	Ask     bool   `short:"a" long:"ask" help:"Ask for a password to set instead of generating one. This is the most secure option."`
}

// Run add
func (cmd *CmdUserAdd) Run(in []string) error {
	if cmd.Help || cmd.Name == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	if cmd.Length < 12 {
		cmd.Length = 12
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
	u, err := db.AddUser(cmd.Name, pw, cmd.Email, cmd.First, cmd.Last, "{}", "{}", cmd.Cost)
	if err != nil {
		return err
	}

	if cmd.Dovecot {
		if cmd.Rounds < 10000 {
			cmd.Rounds = 10000
		}

		u.SetDovecotPassword(pw, cmd.Rounds)
		err = db.SaveUser(u)
		if err != nil {
			return err
		}
	}

	if db.GetFlag("mailmode") {
		err = db.SetAlias(cmd.Name, cmd.Name)
		if err != nil {
			return err
		}
	}

	if !cmd.Ask {
		m("Added user %d with password %s%s%s", u.ID, ansi.Green, pw, ansi.Normal)
	}
	return nil
}
