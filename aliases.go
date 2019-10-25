package anthropoi

// Alias object.
type Alias struct {
	Alias  string `json:"alias"`
	Target string `json:"target"`
}

// Aliases container.
type Aliases struct {
	List []Alias `json:"aliases"`
}

// SetAlias creates or updates a new alias pointing to an existing target address (which may itself be an alias).
func (db *DBM) SetAlias(alias, target string) error {
	q := `INSERT INTO public.aliases (alias,target) 
	VALUES ($1, $2)
	ON CONFLICT (alias) DO UPDATE SET target=$2;`
	_, err := db.Exec(q, alias, target)
	return err
}

// GetAlias returns the target for an alias.
func (db *DBM) GetAlias(alias string) (string, error) {
	var t string
	err := db.QueryRow("SELECT target FROM public.aliases WHERE alias=$1;", alias).Scan(&t)
	if err != nil {
		return "", err
	}

	return t, nil
}

// SearchAliases finds aliases or targets containing the match string. Leave blank to list everything.
func (db *DBM) SearchAliases(match string) (*Aliases, error) {
	rows, err := db.Query("SELECT alias,target FROM public.aliases WHERE alias||target LIKE '%'||$1||'%' ORDER BY target;", match)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var aliases Aliases
	for rows.Next() {
		var a Alias
		err = rows.Scan(&a.Alias, &a.Target)
		if err != nil {
			return nil, err
		}

		aliases.List = append(aliases.List, a)
	}

	return &aliases, nil
}

// RemoveAlias deletes an alias.
func (db *DBM) RemoveAlias(alias string) error {
	_, err := db.Exec("DELETE FROM public.aliases WHERE alias=$1;", alias)
	return err
}

// RemoveAliases deletes all aliases with the same target.
func (db *DBM) RemoveAliases(target string) error {
	_, err := db.Exec("DELETE FROM public.aliases WHERE target=$1;", target)
	return err
}
