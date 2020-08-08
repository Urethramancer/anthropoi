package main

import (
	"fmt"
	"net/http"
)

const (
	errUserPassword       = "User and/or password unknown."
	errorInvalidToken     = "Invalid token. Please authenticate."
	errorBadPassword      = "Bad or easily guessable password."
	errorWrongPassword    = "Password failure. Are you logged in with a different account than you think?"
	errorPasswordMismatch = "Both replacement passwords must match."
)

func apierror(w http.ResponseWriter, msg string) {
	s := fmt.Sprintf("{\"message\":\"%s\",\"ok\":false}", msg)
	w.Write([]byte(s))
}
