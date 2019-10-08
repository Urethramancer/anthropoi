package profiles

import (
	"database/sql"
	"strings"
)

// Profile for a site or business.
// One or more of these exist per user account.
type Profile struct {
	// ID of profile in the database.
	ID int64
	// User ID this profile belongs to.
	User int64
	// Domain this profile is for. This may be an actual Internet
	// domain or the display name for a site/business.
	Domain string
	// Groups the user belongs to on this site, with permissions.
	Groups map[string]string
	// Data for the site. Usually a custom JSON structure.
	Data string
}

// InitTables creates tables and triggers for profiles.
func InitTables(db *sql.DB) error {
	_, err := db.Exec(profileTable)
	if err != nil {
		return err
	}

	return nil
}

// New creates an initialised profile structure.
func New(id, user int64, domain, data string) *Profile {
	p := &Profile{
		ID:     id,
		User:   user,
		Domain: domain,
		Data:   data,
		Groups: make(map[string]string),
	}
	return p
}

// SetGroup adds or creates a group permissions entry to the profile.
func (p *Profile) SetGroup(group string, permissions []string) {
	perm := strings.Join(permissions, ",")
	p.Groups[group] = perm
}

// RemoveGroup if it exists in the profile.
func (p *Profile) RemoveGroup(group string) {
	_, ok := p.Groups[group]
	if !ok {
		return
	}

	delete(p.Groups, group)
}
