package anthropoi

// alias character varying(200) COLLATE pg_catalog."default" NOT NULL,
// target character varying(200) COLLATE pg_catalog."default" NOT NULL,

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
	err := db.QueryRow("SELECT value FROM public.variables WHERE key=$1;", alias).Scan(&t)
	if err != nil {
		return "", err
	}

	return t, nil
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
