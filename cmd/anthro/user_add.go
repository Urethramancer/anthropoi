package main

import (
	"errors"

	"github.com/Urethramancer/ansi"
	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/opt"
)

// CmdUserAdd options.
type CmdUserAdd struct {
	opt.DefaultHelp
	Name  string `placeholder:"USERNAME" help:"Name of user to add." opt:"required"`
	Email string `short:"e" long:"email" help:"Optional e-mail"`
	First string `short:"f" long:"firstname" help:"Optional first name (the one displayed first - may be family name for some regions)."`
	Last  string `short:"l" long:"lastname" help:"Optional last name."`
	Cost  int    `short:"c" long:"cost" help:"Cost of hashing algorithm. Tweak this to make it at least 500ms to hash on the server." default:"12"`
}

// Run add
func (cmd *CmdUserAdd) Run(in []string) error {
	if cmd.Help || cmd.Name == "" {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect()
	if err != nil {
		return err
	}

	defer db.Close()
	pw := anthropoi.GenString(14)
	u, err := db.AddUser(cmd.Name, pw, cmd.Email, cmd.First, cmd.Last, "{}", "{}", cmd.Cost)
	if err != nil {
		return err
	}

	m("Added user %d with password %s%s%s", u.ID, ansi.Green, pw, ansi.Normal)
	return nil
}
