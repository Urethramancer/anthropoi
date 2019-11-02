package main

import (
	"encoding/json"
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
	t := as.getToken(msg.Token)
	if t == nil {
		apierror(w, errorInvalidToken, 403)
		return
	}

	if !t.User.AcceptablePassword(msg.Password) {
		apierror(w, errorBadPassword, 406)
		return
	}

	a := t.User.SplitPasswordElements()
	if len(a) == 4 {
		err := t.User.SetPassword(msg.Password, 0)
		if err != nil {
			apierror(w, err.Error(), 500)
			return
		}
	} else {
		t.User.SetDovecotPassword(msg.Password, 0)
	}

	err := as.db.SaveUser(t.User)
	if err != nil {
		apierror(w, err.Error(), 500)
		return
	}

	reply := StatusReply{}
	reply.Message = "Password changed."
	data, err := json.Marshal(reply)
	if err != nil {
		apierror(w, err.Error(), 500)
		return
	}

	w.Write([]byte(data))
	as.invalidateToken(msg.Token)
	as.L("Password for %s changed by %s", t.User.Username, r.RemoteAddr)
}
