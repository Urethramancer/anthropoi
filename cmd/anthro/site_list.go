package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Urethramancer/signor/stringer"

	"github.com/Urethramancer/signor/opt"
)

// CmdSiteList options.
type CmdSiteList struct {
	opt.DefaultHelp
	Match string `placeholder:"KEYWORD" help:"Find sites containing keyword. Leave blank to list all."`
	JSON  bool   `short:"j" long:"json" help:"Output in JSON format."`
}

// Run list
func (cmd *CmdSiteList) Run(in []string) error {
	if cmd.Help {
		return errors.New(opt.ErrorUsage)
	}

	db, err := connect(name)
	if err != nil {
		return err
	}

	defer db.Close()
	sites, err := db.SearchSites(cmd.Match)
	if err != nil {
		return err
	}

	buf := stringer.New()
	buf.WriteStrings("ID\tDomain\n")
	for _, site := range sites.List {
		buf.WriteI(site.ID, "\t", site.Name, "\n")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, buf.String())
	w.Flush()
	return nil
}
