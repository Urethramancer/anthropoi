package main

import (
	"fmt"
	"net/http"
)

const (
	errUserPassword     = "User and/or password unknown."
	errorInvalidToken   = "Invalid token. Please authenticate."
	errorBadPassword    = "Bad or easily guessable password."
	errorWrongPassword  = "Password failure. Are you logged in with a different account than you think?"
	errPasswordMismatch = "Both replacement passwords must match."
	errMissingEmail     = "Missing e-mail."
	errSettingRecEmail  = "Couldn't set recovery e-mail."
)

func apierror(w http.ResponseWriter, msg string) {
	s := fmt.Sprintf("{\"message\":\"%s\",\"ok\":false}", msg)
	w.Write([]byte(s))
}
