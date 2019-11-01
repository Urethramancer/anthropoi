package main

import (
	"encoding/json"
	"net/http"
)

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
