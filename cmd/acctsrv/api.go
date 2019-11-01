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

// Get details, update details.
func (as *AccountServer) user(w http.ResponseWriter, r *http.Request) {

}

type password struct {
	// Token for the user who's changing the password.
	Token string `json:"token"`
	// Password to set.
	Password string `json:"password"`
}

func (as *AccountServer) password(w http.ResponseWriter, r *http.Request) {
	var msg password
	json.NewDecoder(r.Body).Decode(&msg)
	as.L("%+v", msg)
}
