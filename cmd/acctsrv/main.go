package main

import (
	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/env"
)

func main() {
	as := NewAccountServer(
		env.Get("DB_HOST", "localhost"),
		env.Get("DB_NAME", "accounts"),
		env.Get("DB_USERNAME", "postgres"),
		env.Get("DB_PASSWORD", "postgres"),
	)

	as.Start()
	<-daemon.BreakChannel()
	as.Stop()

}

func NewAccountServer(host, database, username, password string) *AccountServer {
	as := AccountServer{
		dbhost:     host,
		dbname:     database,
		dbusername: username,
		dbpassword: password,
	}

	return &as
}

type AccountServer struct {
	dbhost     string
	dbname     string
	dbusername string
	dbpassword string
}

func (as *AccountServer) Start() {

}

func (as *AccountServer) Stop() {

}
