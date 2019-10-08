package users

import (
	"database/sql"
	"sync"
)

type UserManager struct {
	sync.Mutex
	*sql.DB

	users map[string]*User
}

// NewManager initialises a UserManager.
func NewManager(db *sql.DB) *UserManager {
	um := &UserManager{
		users: make(map[string]*User),
	}
	um.DB = db
	return um
}

// Add a new account.
func (um *UserManager) Add(username, password, email, data string, cost int) (*User, error) {
	u, err := New(username, password, email, data, cost)
	if err != nil {
		return nil, err
	}

	um.users[username] = u
	return u, nil
}
