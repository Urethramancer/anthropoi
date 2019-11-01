package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Urethramancer/anthropoi"
)

func (as *AccountServer) authenticate(w http.ResponseWriter, r *http.Request) {
	var msg AuthMsg
	json.NewDecoder(r.Body).Decode(&msg)
	reply := StatusReply{}

	u, err := as.db.GetUserByName(msg.Username)
	if err != nil {
		apierror(w, errUserPassword, 403)
		return
	}
	println("2")

	as.L("%s: %s", u.Usermame, u.Password)
	as.L("%s", anthropoi.GenerateDovecotPassword(msg.Password, u.Salt, 50000))
	if !u.CompareDovecotHashAndPassword(msg.Password) {
		apierror(w, errUserPassword, 403)
		return
	}

	h := sha256.New()
	h.Write([]byte(anthropoi.GenString(16)))
	h.Write([]byte(time.Now().String()))
	reply.Message = hex.EncodeToString(h.Sum(nil))

	data, err := json.Marshal(reply)
	if err != nil {
		apierror(w, err.Error(), 500)
		return
	}

	t := Token{
		User:      u,
		Timestamp: time.Now(),
	}
	as.tokens[reply.Message] = &t
	w.Write([]byte(data))
}
