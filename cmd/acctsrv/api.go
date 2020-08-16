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

// UserResponse is a stripped-down version of the internal User structure.
type UserResponse struct {
	// ID of user in the database.
	ID int64 `json:"id"`
	// Username to log in with.
	Username string `json:"username"`
	// Email to verify account or reset password.
	Email string `json:"email"`
	// Created timestamp.
	Created time.Time `json:"created"`
	// First name of user (optional).
	First string `json:"first"`
	// Last name of user (optional).
	Last string `json:"last"`
	// Locked accounts can't log in.
	Locked bool `json:"locked"`
	// Admin for the whole system if true.
	Admin bool `json:"admin"`
	// OK is true.
	OK bool `json:"ok"`
}

func notfound(w http.ResponseWriter, r *http.Request) {
	apierror(w, "Unknown endpoint.")
}

func preflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	w.Header().Set("Access-Control-Max-Age", "86400")
	http.Error(w, "", 204)
}

// Get user details.
func (as *AccountServer) getuser(w http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value("req").(RequestMsg)
	t := as.getToken(msg.Token)
	if t == nil {
		apierror(w, errorInvalidToken)
		return
	}

	res := UserResponse{}
	res.ID = t.User.ID
	res.Username = t.User.Username
	res.Email = t.User.Email
	res.Created = t.User.Created
	res.First = t.User.First
	res.Last = t.User.Last
	res.Locked = t.User.Locked
	res.Admin = t.User.Admin
	res.OK = true
	data, err := json.Marshal(res)
	if err != nil {
		apierror(w, err.Error())
		return
	}

	w.Write([]byte(data))
}

func (as *AccountServer) password(w http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value("req").(RequestMsg)
	t := as.getToken(msg.Token)
	if t == nil {
		apierror(w, errorInvalidToken)
		return
	}

	if msg.Password != msg.PasswordAgain {
		apierror(w, errorPasswordMismatch)
		return
	}

	if !t.User.CompareDovecotHashAndPassword(msg.CurrentPassword) {
		apierror(w, errorWrongPassword)
		return
	}

	if !t.User.AcceptablePassword(msg.Password) {
		apierror(w, errorBadPassword)
		return
	}

	a := t.User.SplitPasswordElements()
	if len(a) == 4 {
		err := t.User.SetPassword(msg.Password, 0)
		if err != nil {
			apierror(w, err.Error())
			return
		}
	} else {
		t.User.SetDovecotPassword(msg.Password, 0)
	}

	err := as.db.SaveUser(t.User)
	if err != nil {
		apierror(w, err.Error())
		return
	}

	reply := StatusReply{}
	reply.Message = "Password changed."
	reply.OK = true
	data, err := json.Marshal(reply)
	if err != nil {
		apierror(w, err.Error())
		return
	}

	w.Write([]byte(data))
	as.invalidateToken(msg.Token)
	as.L("Password for %s changed by %s", t.User.Username, r.RemoteAddr)
}

func (as *AccountServer) aliases(w http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value("req").(RequestMsg)
	t := as.getToken(msg.Token)
	if t == nil {
		apierror(w, errorInvalidToken)
		return
	}

	a, err := as.db.GetAliasesForUser(t.User)
	if err != nil {
		apierror(w, err.Error())
		return
	}

	if a.List[0].Alias == a.List[0].Target {
		a.List = a.List[1:]
	}
	data, err := json.Marshal(a)
	if err != nil {
		apierror(w, err.Error())
		return
	}

	w.Write([]byte(data))
}
