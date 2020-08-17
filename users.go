package anthropoi

import (
	"fmt"
	"time"
)

// User account structure holds basic login and personal information.
type User struct {
	/*
	 * Required bits
	 */

	// ID of user in the database.
	ID int64 `json:"id"`
	// Username to log in with.
	Username string `json:"username"`
	// Password for user account.
	Password string `json:"password"`
	// Salt for the password.
	Salt string `json:"salt"`
	// Email to verify account or reset password.
	Email string `json:"email"`
	// Created timestamp.
	Created time.Time `json:"created"`

	/*
	 * Optional bits
	 */

	// First name of user (optional).
	First string `json:"first"`
	// Last name of user (optional).
	Last string `json:"last"`
	// Data for the account. JSON field for all the customising you need.
	Data string `json:"data"`
	// Tokens is meant to store any authentication tokens required for external sites.
	Tokens string `json:"token"`

	// Sites the user is a member of.
	Sites []string

	// Locked accounts can't log in.
	Locked bool `json:"locked"`
	// Admin for the whole system if true.
	Admin bool `json:"admin"`
}

// Users container.
type Users struct {
	List []*User `json:"users"`
}

const (
	pre_bcrypt      = "$2a$"
	pre_sha512crypt = "{SHA512-CRYPT}"
)

// AddUser creates a new User. This may fail.
func (db *DBM) AddUser(username, password, email, first, last, data, tokens string, cost int) (*User, error) {
	u := &User{
		Username: username,
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
	err = st.QueryRow(u.Username, u.Password, u.Salt, u.Email, u.First, u.Last, u.Data, u.Tokens).Scan(&u.ID)
	return u, err
}

// UpdateUser saves an existing user by ID.
func (db *DBM) SaveUser(u *User) error {
	if u.Data == "" {
		u.Data = "{}"
	}

	if u.Tokens == "" {
		u.Tokens = "{}"
	}

	q := `UPDATE public.users SET username=$1,password=$2,salt=$3,email=$4,locked=$5,first=$6,last=$7,data=$8,tokens=$9,admin=$10 WHERE id=$11;`
	_, err := db.Exec(q, u.Username, u.Password, u.Salt, u.Email, u.Locked, u.First, u.Last, u.Data, u.Tokens, u.Admin, u.ID)
	if err != nil {
		fmt.Printf("WTF? %s\n", err.Error())
		return err
	}

	return nil
}

// GetUser returns a User based on an ID.
func (db *DBM) GetUser(id int64) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens,admin FROM public.users WHERE id=$1 LIMIT 1", id).Scan(
		&u.ID, &u.Username, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens, &u.Admin)
	if err != nil {
		return nil, err
	}

	return &u, db.GetSitesForUser(&u)
}

// GetUserByName for when you don't have an ID.
func (db *DBM) GetUserByName(name string) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens,admin FROM public.users WHERE username=$1 LIMIT 1", name).Scan(
		&u.ID, &u.Username, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens, &u.Admin)
	if err != nil {
		return nil, err
	}

	return &u, db.GetSitesForUser(&u)
}

// GetSitesForUser fills the Sites field in the User struct.
func (db *DBM) GetSitesForUser(u *User) error {
	q := `SELECT name FROM public.users u
	INNER JOIN membership m ON u.id=m.userid
	INNER JOIN sites s ON m.siteid=s.id WHERE u.id=$1;`
	rows, err := db.Query(q, u.ID)
	if err != nil {
		return err
	}

	defer rows.Close()
	u.Sites = []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		u.Sites = append(u.Sites, name)
	}
	return nil
}

// RemoveUser by ID.
func (db *DBM) RemoveUser(id int64) error {
	_, err := db.Exec("DELETE FROM public.users WHERE id=$1;", id)
	return err
}

// RemoveUserByName for when that's needed.
func (db *DBM) RemoveUserByName(name string) error {
	_, err := db.Exec("DELETE FROM public.users WHERE username=$1;", name)
	return err
}

// GetUsers retrieves users, sorted by ID, optionally containing a keyword.
func (db *DBM) GetUsers(match string) (*Users, error) {
	rows, err := db.Query("SELECT id,username,email,created,locked,first,last,admin FROM public.users WHERE username LIKE '%'||$1||'%' ORDER BY id;", match)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users Users
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Admin)
		if err != nil {
			return nil, err
		}

		err = db.GetSitesForUser(&u)
		if err != nil {
			return nil, err
		}

		users.List = append(users.List, &u)
	}

	return &users, nil
}
