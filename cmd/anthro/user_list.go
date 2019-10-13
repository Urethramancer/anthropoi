package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
)

// CmdUserList options.
type CmdUserList struct {
	opt.DefaultHelp
}

// Run list
func (cmd *CmdUserList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect()
	if err != nil {
		return err
	}

	defer db.Close()
	list, err := db.GetUsers(0)
	if err != nil {
		return err
	}

	var out stringer.Stringer
	out.WriteStrings("ID\tUsername\tName\tE-mail\tCreated\tActive\n")
	for _, u := range list {
		if u.First == "" && u.Last == "" {
			u.First = "<unset>"
		}
		if u.Email == "" {
			u.Email = "<unset>"
		}
		out.WriteI(
			fmt.Sprintf("%d\t", u.ID),
			u.Usermame, "\t",
			u.First, " ", u.Last, "\t",
			u.Email, "\t",
			u.Created.String(), "\t",
			!u.Locked, "\t",
			"\n",
		)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, out.String())
	w.Flush()
	return nil
}
