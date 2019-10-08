package groups

import (
	"database/sql"
	"sync"
)

// GroupManager handles the lifetime of Group structures.
type GroupManager struct {
	sync.Mutex
	*sql.DB

	// Groups in this manager.
	Groups map[string]*Group

	// next ID for the database
	next int64
}

// NewManager initialises a GroupManager.
func NewManager(db *sql.DB) *GroupManager {
	gm := &GroupManager{
		Groups: make(map[string]*Group),
		next:   1, // TODO: Grab from DB
	}
	gm.DB = db
	return gm
}

// InitTable initialises the Group table in the database.
func (gm *GroupManager) InitTable() error {
	return nil
}

// Add a group and permissions.
func (gm *GroupManager) Add(name string, permissions []string) {
	gm.Lock()
	defer gm.Unlock()
	g := New(gm.next, name, permissions)
	gm.Groups[name] = g
	gm.next++
}

// Remove a group.
func (gm *GroupManager) Remove(name string) {
	gm.Lock()
	defer gm.Unlock()
	_, ok := gm.Groups[name]
	if !ok {
		return
	}

	delete(gm.Groups, name)
}

// AddPermissions adds the supplied permissions to a group.
func (gm *GroupManager) AddPermissions(group string, p ...string) {
	gm.Lock()
	defer gm.Unlock()
	g, ok := gm.Groups[group]
	if !ok {
		return
	}

	g.AddPermissions(p...)
}

// RemovePermissions removes the supplied permissions from a group.
func (gm *GroupManager) RemovePermissions(group string, p ...string) {
	gm.Lock()
	defer gm.Unlock()
	g, ok := gm.Groups[group]
	if !ok {
		return
	}

	for _, perm := range p {
		g.RemovePermission(perm)
	}
}
