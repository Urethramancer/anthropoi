package anthropoi

import (
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
func (db *DBM) AddUser(username, password, email, data string, cost int) (*User, error) {
	u := &User{
		Usermame: username,
		Email:    email,
		Salt:     genString(32),
		Data:     data,
	}
	err := u.SetPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUser returns a User based on an ID.
func (db *DBM) GetUser(id int64) (*User, error) {
	var u User
	err := db.QueryRow("SELECT * FROM public.users WHERE id='?'", id).Scan(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByName for when you don't have an id.
func (db *DBM) GetUserByName(name string) (*User, error) {
	var u User
	err := db.QueryRow("", name).Scan(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// SaveUser via upsert.
func (db *DBM) SaveUser(u *User) (int64, error) {
	res, err := db.Exec("INSERT INTO public.users (id,username,password,salt,email,locked,first,last,data,tokens) VALUES (?,?,?,?,?,?,?,?,?,?) ON CONFLICT ON CONSTRAINT key_users_pkey DO UPDATE SET username=EXCLUDED.username,password=EXCLUDED.password,salt=EXCLUDED.salt,email=EXCLUDED.email,locked=EXCLUDED.locked,first=EXCLUDED.first,last=EXCLUDED.last,data=EXCLUDED.data,tokens=EXCLUDED.tokens RETURNING id;",
		u.ID, u.Usermame, u.Password, u.Salt, u.Email, u.Locked, u.First, u.Last, u.Data, u.Tokens)
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
