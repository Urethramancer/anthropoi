package main

import (
	"net/http"
	"time"

	"github.com/Urethramancer/anthropoi"
)

// Token keeps track of when it was last used, and who it's for.
type Token struct {
	User      *anthropoi.User
	Timestamp time.Time
}

func notfound(w http.ResponseWriter, r *http.Request) {
	apierror(w, "Unknown endpoint.", 404)
}

// Get details, update details.
func (as *AccountServer) user(w http.ResponseWriter, r *http.Request) {

}

func (as *AccountServer) setPassword(w http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value("req").(RequestMsg)
	as.L("%s, %s", msg.Token, msg.Password)
	t := as.getToken(msg.Token)
	if t == nil {
		apierror(w, errorInvalidToken, 403)
		return
	}

	as.L("Setting password for %s from %s", t.User.Username, r.RemoteAddr)
}
