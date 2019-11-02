package main

// RequestMsg is a combined structure for all calls.
type RequestMsg struct {
	// Token is returned from authentication.
	Token string `json:"token,omitempty"`
	// Username is required when authenticating.
	Username string `json:"username,omitempty"`
	// Password to authenticate or set.
	Password string `json:"password,omitempty"`
}

// StatusReply is returned from all calls.
type StatusReply struct {
	// Message string.
	Message string `json:"message"`
	// Code is 0 if all went well. If this was embedded in another struct, there might be other data.
	Code int `json:"code"`
}
