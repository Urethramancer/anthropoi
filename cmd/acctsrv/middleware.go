package main

import (
	"context"
	"net/http"
)

// AuthMsg for authentication.
type AuthMsg struct {
	// Username is required.
	Username string `json:"username"`
	// Password is required.
	Password string `json:"password"`
}

// StatusReply is returned from all calls.
type StatusReply struct {
	// Message string.
	Message string `json:"message"`
	// OK is true if all went well. If this was embedded in another struct, there will be other data.
	OK bool `json:"ok"`
}

// Check security token.
func check_access(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "Authentication", "moo"))
		w.Write([]byte("auth"))
		// http.Error(w, "Unknown token", http.StatusForbidden)
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
