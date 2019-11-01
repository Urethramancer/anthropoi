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

	if !u.CompareDovecotHashAndPassword(msg.Password) {
		as.L("User %s failed to authenticate from %s", u.Username, r.RemoteAddr)
		apierror(w, errUserPassword, 403)
		return
	}

	as.L("User %s authenticated from %s", u.Username, r.RemoteAddr)
	reply.Message = as.createToken(u)
	data, err := json.Marshal(reply)
	if err != nil {
		apierror(w, err.Error(), 500)
		return
	}

	w.Write([]byte(data))
}

// Create or get an active token.
func (as *AccountServer) createToken(u *anthropoi.User) string {
	hash, ok := as.hashes[u.Username]
	if ok {
		past := time.Now().Add(-time.Minute * 30)
		if !as.tokens[hash].Timestamp.Before(past) {
			as.tokens[hash].Timestamp = time.Now()
			return hash
		} else {
			delete(as.hashes, u.Username)
			delete(as.tokens, hash)
		}
	}

	h := sha256.New()
	h.Write([]byte(anthropoi.GenString(16)))
	h.Write([]byte(time.Now().String()))
	hash = hex.EncodeToString(h.Sum(nil))

	t := Token{
		User:      u,
		Timestamp: time.Now(),
	}
	as.hashes[u.Username] = hash
	as.tokens[hash] = &t
	return hash
}
