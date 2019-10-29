package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/env"
	"github.com/Urethramancer/signor/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	as := NewAccountServer(
		env.Get("DB_HOST", "localhost"),
		env.Get("DB_NAME", "accounts"),
		env.Get("DB_USERNAME", "postgres"),
		env.Get("DB_PASSWORD", "postgres"),
		env.Get("WEB_NOST", "127.0.0.1"),
		env.Get("WEB_PORT", "8000"),
	)

	as.Start()
	<-daemon.BreakChannel()
	as.Stop()
}

// NewAccountServer sets up routes and returns a server ready to start serving REST end points.
func NewAccountServer(dbhost, database, username, password, host, port string) *AccountServer {
	as := AccountServer{
		dbhost:     dbhost,
		dbname:     database,
		dbusername: username,
		dbpassword: password,
		host:       host,
		port:       port,
	}

	as.L = log.Default.TMsg
	as.E = log.Default.TErr

	as.IdleTimeout = time.Second * 30
	as.ReadTimeout = time.Second * 10
	as.WriteTimeout = time.Second * 10

	as.api = chi.NewRouter()
	as.api.Use(middleware.NoCache)
	as.api.Use(middleware.RealIP)
	as.api.Use(middleware.RequestID)
	as.api.Use(middleware.Timeout(time.Second * 10))

	as.api.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("1"))
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(check_access)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("/user"))
			})
		})

	})

	return &as
}

type AccountServer struct {
	sync.RWMutex
	sync.WaitGroup
	log.LogShortcuts
	http.Server
	dbhost     string
	dbname     string
	dbusername string
	dbpassword string
	host       string
	port       string

	api *chi.Mux
}

// Start the server.
func (as *AccountServer) Start() {
	as.Lock()
	defer as.Unlock()
	addr := net.JoinHostPort(as.host, as.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		as.E("Listener error: %s", err.Error())
		os.Exit(2)
	}

	as.Add(1)
	as.L("Starting web server on http://%s", addr)
	go func() {
		as.Handler = as.api
		err = as.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			as.E("Error running server: %s", err.Error())
			os.Exit(2)
		}
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
}

func check_access(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "Authentication", "moo"))
		w.Write([]byte("auth"))
		// http.Error(w, "Unknown token", http.StatusForbidden)
	}
	return http.HandlerFunc(fn)
}