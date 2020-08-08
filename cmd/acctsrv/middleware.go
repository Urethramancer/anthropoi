package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (as *AccountServer) decode_request(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var msg RequestMsg
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			apierror(w, err.Error())
			as.E("Error decoding JSON for request from %s: %s", r.RemoteAddr, err.Error())
			return
		}

		r = r.WithContext(context.WithValue(ctx, "req", msg))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Check security token.
func (as *AccountServer) check_access(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		msg := r.Context().Value("req").(RequestMsg)
		t := as.getToken(msg.Token)
		if t == nil {
			apierror(w, errorInvalidToken)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func addJSONHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func addCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
