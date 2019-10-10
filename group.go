package anthropoi

import (
	"sort"
)

// Group structures contain the ID, name and a list of possible permissions.
type Group struct {
	// ID of group in the database.
	ID int64
	// Name of group.
	Name string
	// Permissions which can be set in this group.
	Permissions []string
}

// New creates an initialised group structure.
func New(id int64, name string, permissions []string) *Group {
	g := &Group{
		ID:          id,
		Name:        name,
		Permissions: permissions,
	}
	return g
}

// AddPermissions and sort.
func (g *Group) AddPermissions(p ...string) {
	g.Permissions = append(g.Permissions, p...)
	sort.Strings(g.Permissions)
}

// RemovePermission from group.
func (g *Group) RemovePermission(p string) {
	var list []string
	for _, perm := range g.Permissions {
		if perm != p {
			list = append(list, perm)
		}
	}
	g.Permissions = list
}
