package anthropoi

import (
	"time"
)

// Site or domain.
type Site struct {
	ID      int64            `json:"id"`
	Name    string           `json:"name"`
	Created time.Time        `json:"created"`
	Groups  map[string]Group `json:"groups"`
}

// Sites container.
type Sites struct {
	List []*Site `json:"sites"`
}

// Group for a site.
type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SearchSites finds sites containing the match string. Leave blank to list everything.
func (db *DBM) SearchSites(match string) (*Sites, error) {
	rows, err := db.Query("SELECT id,name,created FROM public.sites WHERE name LIKE '%'||$1||'%' ORDER BY id;", match)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sites Sites
	for rows.Next() {
		var s Site
		err = rows.Scan(&s.ID, &s.Name, &s.Created)
		if err != nil {
			return nil, err
		}

		sites.List = append(sites.List, &s)
	}

	return &sites, nil
}

// AddSite to enable users being associated.
func (db *DBM) AddSite(name string) (int64, error) {
	st, err := db.Prepare("INSERT INTO public.sites (name) VALUES($1) RETURNING ID;")
	if err != nil {
		return 0, err
	}

	defer st.Close()
	var id int64
	err = st.QueryRow(name).Scan(&id)
	return id, err
}

// RemoveSite by ID.
func (db *DBM) RemoveSite(id int64) error {
	_, err := db.Exec("DELETE FROM public.sites WHERE id=$1;", id)
	return err
}

// RemoveSiteByName for when that's more convenient.
func (db *DBM) RemoveSiteByName(name string) error {
	_, err := db.Exec("DELETE FROM public.sites WHERE name=$1;", name)
	return err
}
