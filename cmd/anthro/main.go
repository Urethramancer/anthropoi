package main

import (
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	_ "github.com/lib/pq"
)

// Options holds all the tool commands.
var Options struct {
	opt.DefaultHelp
	Stamp bool     `short:"t" long:"timestamp" help:"Timestamp all output."`
	Init  CmdInit  `command:"init" help:"Initialise database and tables."`
	User  CmdUser  `command:"user" help:"User management." aliases:"u"`
	Group CmdGroup `command:"group" help:"Group management." aliases:"g"`
}

var m func(string, ...interface{})
var e func(string, ...interface{})

func main() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	if Options.Stamp {
		m = log.Default.TMsg
		e = log.Default.TErr
	} else {
		m = log.Default.Msg
		e = log.Default.Err
	}

	err := a.RunCommand(false)
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}
}

func getenv(key, alt string) string {
	s := os.Getenv(key)
	if s == "" {
		return alt
	}

	return s
}
