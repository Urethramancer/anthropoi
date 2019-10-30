package main

import (
	"encoding/json"
	"net/http"
)

func (as *AccountServer) authenticate(w http.ResponseWriter, r *http.Request) {
	var msg AuthMsg
	json.NewDecoder(r.Body).Decode(&msg)
	as.L("%+v", msg)
	reply := StatusReply{}
	reply.Message = "token"
	reply.OK = true
	data, err := json.Marshal(reply)
	if err != nil {
		as.E("Error marshalling: %s", err.Error())
		return
	}
	w.Write([]byte(data))
}
