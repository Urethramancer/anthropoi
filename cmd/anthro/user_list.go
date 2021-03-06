package main

import (
	"encoding/json"
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
	Match string `placeholder:"KEYWORD" help:"Find users containing keyword. Leave blank to list all."`
	JSON  bool   `short:"j" long:"json" help:"Output in JSON format."`
}

// Run list
func (cmd *CmdUserList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	users, err := db.GetUsers(cmd.Match)
	if err != nil {
		return err
	}

	if cmd.JSON {
		j, err := json.MarshalIndent(users, "", "\t")
		if err != nil {
			return err
		}

		m("%s", string(j))
	} else {
		var out stringer.Stringer
		out.WriteStrings("ID\tUsername\tName\tE-mail\tDomains\tCreated\tActive\tAdmin\n")
		for _, u := range users.List {
			if u.First == "" && u.Last == "" {
				u.First = "<unset>"
			}
			if u.Email == "" {
				u.Email = "<unset>"
			}
			out.WriteI(
				fmt.Sprintf("%d\t", u.ID),
				u.Username, "\t",
				u.First, " ", u.Last, "\t",
				u.Email, "\t",
				len(u.Sites), "\t",
				u.Created.String(), "\t",
				!u.Locked, "\t",
				u.Admin, "\t",
				"\n",
			)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprint(w, out.String())
		w.Flush()
	}
	return nil
}
