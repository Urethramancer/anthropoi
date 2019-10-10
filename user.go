package anthropoi

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User account structure holds basic login and personal information.
type User struct {
	/*
	 * Required bits
	 */

	// ID of user in the database.
	ID int64
	// Username to log in with.
	Usermame string
	// Password for user account.
	Password string
	// Salt for the password.
	Salt string
	// Email to verify account or reset password.
	Email string
	// Created timestamp.
	Created time.Time
	// Locked accounts can't log in.
	Locked bool

	/*
	 * Optional bits
	 */

	// Profiles for specific sites.
	Profiles []*Profile

	// First name of user (optional).
	First string
	// Last name of user (optional).
	Last string
	// Data for the account. JSON field for all the customising you need.
	Data string
	// Tokens is meant to store any authentication tokens required for external sites.
	Tokens string
}

// AddUser creates an initialised User structure. This may fail.
func AddUser(username, password, email, data string, cost int) (*User, error) {
	acc := &User{
		Usermame: username,
		Email:    email,
		Salt:     genString(32),
		Data:     data,
	}
	err := acc.SetPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// GetUser returns a User based on an ID.
func GetUser(db *sql.DB, id int64) *User {
	var u User
	// q := "SELECT * FROM public.users WHERE"
	return &u
}

// SetPassword to a new one.
func (u *User) SetPassword(password string, cost int) error {
	s := password + u.Salt
	hash, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword against the account's hash.
func (u *User) CheckPassword(password string) bool {
	if u.Usermame == "" || u.Password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return (err == nil)
}
