package main

// RequestMsg is a combined structure for all calls.
type RequestMsg struct {
	// Token is returned from authentication.
	Token string `json:"token,omitempty"`
	// Username is required when authenticating.
	Username string `json:"username,omitempty"`
	// Password to authenticate or set.
	Password string `json:"password,omitempty"`
	// CurrentPassword is used in password change requests.
	CurrentPassword string `json:"currentpassword,omitempty"`
	// PasswordAgain is used for server-side form validation in change requests.
	PasswordAgain string `json:"passwordagain,omitempty"`
}

// StatusReply is returned from all calls.
type StatusReply struct {
	// Message string.
	Message string `json:"message"`
	// OK is false true if all went well. On false, the message is an error message.
	OK bool `json:"ok"`
}
