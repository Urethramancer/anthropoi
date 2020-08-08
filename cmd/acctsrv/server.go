package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Urethramancer/anthropoi"
	"github.com/Urethramancer/signor/env"
	"github.com/Urethramancer/signor/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

// AccountServer embeds the HTTP server struct.
type AccountServer struct {
	sync.RWMutex
	sync.WaitGroup
	log.LogShortcuts
	http.Server

	dbhost string
	dbport string
	dbname string
	dbuser string
	dbpass string
	db     *anthropoi.DBM

	host string
	port string

	api    *chi.Mux
	hashes map[string]string
	tokens map[string]*Token
}

// NewAccountServer sets up routes and returns a server ready to start serving REST end points.
func NewAccountServer(dbhost, dbport, dbname, dbuser, dbpass, host, port string) *AccountServer {
	as := AccountServer{
		dbhost: dbhost,
		dbport: dbport,
		dbname: dbname,
		dbuser: dbuser,
		dbpass: dbpass,
		host:   host,
		port:   port,
		tokens: make(map[string]*Token),
		hashes: make(map[string]string),
	}

	as.L = log.Default.TMsg
	as.E = log.Default.TErr

	as.IdleTimeout = time.Second * 30
	as.ReadTimeout = time.Second * 10
	as.WriteTimeout = time.Second * 10

	as.api = chi.NewRouter()
	as.api.Use(middleware.NoCache)
	as.api.Use(addCORS)
	as.api.Use(middleware.RealIP)
	as.api.Use(middleware.RequestID)
	as.api.Use(middleware.Timeout(time.Second * 10))
	as.api.NotFound(notfound)

	as.api.Route("/", func(r chi.Router) {
		r.Use(addJSONHeaders)
		r.Use(as.decode_request)
		r.Options("/", preflight)
		r.Post("/auth", as.authenticate)
		r.Post("/password", as.password)
		r.Post("/aliases", as.aliases)
		r.Get("/user", as.user)
	})

	return &as
}

// Start the server.
func (as *AccountServer) Start() {
	as.Lock()
	defer as.Unlock()
	as.db = anthropoi.New(as.dbhost, as.dbport, as.dbuser, as.dbpass, env.Get("DB_SSL", "disable"))
	err := as.db.Connect(as.dbname)
	if err != nil {
		as.E("Couldn't open database: %s", err.Error())
		os.Exit(2)
	}

	addr := net.JoinHostPort(as.host, as.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		as.E("Listener error: %s", err.Error())
		as.db.Close()
		os.Exit(2)
	}

	as.Add(1)
	as.L("Starting web server on http://%s", addr)
	go func() {
		as.Handler = as.api
		err = as.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			as.E("Error running server: %s", err.Error())
			as.db.Close()
			os.Exit(2)
		}
		as.L("%d entries in token map", len(as.tokens))
		as.L("Stopped web server.")
		as.Done()
	}()
}

// Stop the server.
func (as *AccountServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := as.Shutdown(ctx)
	if err != nil {
		as.E("Shutdown error: %s", err.Error())
		os.Exit(2)
	}

	as.Wait()
	as.db.Close()
}
