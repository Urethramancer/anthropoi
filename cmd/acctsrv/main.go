package main

import (
	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/env"
)

func main() {
	as := NewAccountServer(
		env.Get("DB_HOST", "localhost"),
		env.Get("DB_PORT", "5432"),
		env.Get("DB_NAME", "accounts"),
		env.Get("DB_USERNAME", "postgres"),
		env.Get("DB_PASSWORD", ""),
		env.Get("WEB_HOST", "0.0.0.0"),
		env.Get("WEB_PORT", "8000"),
	)

	as.Start()
	if as.db.DatabaseExists("accounts") {
		println("DB OK!")
	}
	err := as.db.Ping()
	if err != nil {
		as.E("DB ping: %s", err.Error())
	}
	<-daemon.BreakChannel()
	as.Stop()
}
