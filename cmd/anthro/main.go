package main

import (
	"os"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/env"
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	_ "github.com/lib/pq"
)

// Options holds all the tool commands.
var Options struct {
	opt.DefaultHelp
	Stamp    bool   `short:"t" long:"timestamp" help:"Timestamp all output."`
	Host     string `short:"H" long:"host" help:"Host to connect to." default:"localhost"`
	Port     string `short:"p" long:"port" help:"Port to connect to." default:"5432"`
	Username string `short:"u" long:"user" help:"User to connect as." default:"postgres"`
	Password string `short:"P" long:"password" help:"Password for that user. Nay be left out if PostgreSQL is configured for other authentication methods."`
	Name     string `short:"n" long:"name" help:"Name of the database to create." default:"accounts"`
	SSL      bool   `short:"s" long:"ssl" help:"Require SSL to connect."`

	Init  CmdInit  `command:"init" help:"Initialise database and tables."`
	Site  CmdSite  `command:"site" help:"Site management." aliases:"s"`
	User  CmdUser  `command:"user" help:"User management." aliases:"u"`
	Alias CmdAlias `command:"alias" help:"Alias management." aliases:"al"`
	Group CmdGroup `command:"group" help:"Group management." aliases:"g"`
}

var m func(string, ...interface{})
var e func(string, ...interface{})

var (
	host     string
	port     string
	username string
	password string
	name     string
	ssl      string
)

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

	host = env.Get("DB_HOST", Options.Host)
	port = env.Get("DB_PORT", "5432")
	username = env.Get("DB_USERNAME", "postgres")
	password = env.Get("DB_PASSWORD", "")
	name = env.Get("DB_NAME", anthropoi.DefaultName)
	if Options.SSL {
		ssl = "enable"
	} else {
		ssl = env.Get("DB_SSL", "disable")
	}

	err := a.RunCommand(false)
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}
}

func connect(dbname string) (*anthropoi.DBM, error) {
	db := anthropoi.New(host, port, username, password, ssl)
	return db, db.Connect(dbname)
}
