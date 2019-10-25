package anthropoi

// Site or domain.
type Site struct {
	ID     int64            `json:"id"`
	Name   string           `json:"name"`
	Groups map[string]Group `json:"groups"`
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
	rows, err := db.Query("SELECT id,name FROM public.sites WHERE name LIKE '%'||$1||'%' ORDER BY id;", match)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sites Sites
	for rows.Next() {
		var s Site
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}

		sites.List = append(sites.List, &s)
	}

	return &sites, nil
}
