package main

import (
	"fmt"
	"net/http"
)

const (
	errUserPassword = "User and/or password unknown."
)

func apierror(w http.ResponseWriter, msg string, code int) {
	s := fmt.Sprintf("{\"message\":\"%s\",\"ok\":false}", msg)
	w.Write([]byte(s))
}
