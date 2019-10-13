package anthropoi

import (
	"database/sql"
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

	if cost < 10 {
		cost = 10
	}

	err := u.SetPassword(password, cost)
	if err != nil {
		return nil, err
	}

	q := "INSERT INTO public.users (username,password,salt,email,first,last,data,tokens) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id;"
	st, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}

	defer st.Close()
	err = st.QueryRow(u.Usermame, u.Password, u.Salt, u.Email, u.First, u.Last, u.Data, u.Tokens).Scan(&u.ID)
	if err != nil {
		return nil, err
	}

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
	_, err := db.Exec("DELETE FROM public.users WHERE id=$1;", id)
	return err
}

// DeleteUserByName for when that's needed.
func (db *DBM) DeleteUserByName(name string) error {
	_, err := db.Exec("DELETE FROM public.users WHERE username=$1;", name)
	return err
}

// GetUsers retrieves all users, up to a limit, sorted by ID.
func (db *DBM) GetUsers(limit int64) ([]*User, error) {
	q := "SELECT id,username,email,created,locked,first,last FROM public.users"
	if limit > 0 {
		q += " LIMIT $1"
	}

	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = db.Query(q, limit)
	} else {
		rows, err = db.Query(q)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Usermame, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last)
		if err != nil {
			return nil, err
		}

		list = append(list, &u)
	}
	return list, nil
}

// SetPassword generates a new salt and sets the password.
func (u *User) SetPassword(password string, cost int) error {
	if cost < 10 {
		cost = 10
	}

	u.Salt = GenString(32)
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
