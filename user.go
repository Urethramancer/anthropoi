package anthropoi

import (
	"fmt"
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

// AddUser creates a new User. This may fail.
func (db *DBM) AddUser(username, password, email, first, last, data, tokens string, cost int) (*User, error) {
	u := &User{
		Usermame: username,
		Email:    email,
		First:    first,
		Last:     last,
		Data:     data,
		Tokens:   tokens,
	}
	err := u.SetPassword(password, cost)
	if err != nil {
		return nil, err
	}

	var id int64
	q := "INSERT INTO public.users"
	err = db.QueryRow(q).Scan(&id)
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

// UpdateUser saves an existing user by ID.
func (db *DBM) SaveUser(u *User) error {
	if u.Data == "" {
		u.Data = "{}"
	}

	if u.Tokens == "" {
		u.Tokens = "{}"
	}

	q := `UPDATE public.users SET username=$1,password=$2,salt=$3,email=$4,locked=$5,first=$6,last=$7,data=$8,tokens=$9 WHERE id=$10;`
	_, err := db.Exec(q, u.Usermame, u.Password, u.Salt, u.Email, u.Locked, u.First, u.Last, u.Data, u.Tokens, u.ID)
	if err != nil {
		fmt.Printf("WTF? %s\n", err.Error())
		return err
	}

	return nil
}

// GetUser returns a User based on an ID.
func (db *DBM) GetUser(id int64) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens FROM public.users WHERE id=$1 LIMIT 1", id).Scan(
		&u.ID, &u.Usermame, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByName for when you don't have an id.
func (db *DBM) GetUserByName(name string) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens FROM public.users WHERE username=$1 LIMIT 1", name).Scan(
		&u.ID, &u.Usermame, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// DeleteUser by ID.
func (db *DBM) DeleteUser(id int64) error {
	_, err := db.Exec("DELETE FROM public.users WHERE id=$1 LIMIT 1;", id)
	return err
}

// DeleteUserByName for when that's needed.
func (db *DBM) DeleteUserByName(name string) error {
	_, err := db.Exec("DELETE FROM public.users WHERE username=$1 LIMIT 1;", name)
	return err
}

// SetPassword generates a new salt and sets the password.
func (u *User) SetPassword(password string, cost int) error {
	u.Salt = genString(32)
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
