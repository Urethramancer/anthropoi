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
		env.Get("DB_PASSWORD", "postgres"),
		env.Get("WEB_HOST", "127.0.0.1"),
		env.Get("WEB_PORT", "8000"),
	)

	as.Start()
	<-daemon.BreakChannel()
	as.Stop()
}
