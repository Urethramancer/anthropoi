package anthropoi

// SetVar sets a string in the variables table.
func (db *DBM) SetVar(key, value string) error {
	q := `INSERT INTO public.variables (key,value) 
	VALUES ($1, $2)
	ON CONFLICT (key) DO UPDATE SET value = $2;`
	_, err := db.Exec(q, key, value)
	return err
}

// GetVar gets a string from the variables table.
func (db *DBM) GetVar(key string) (string, error) {
	var s string
	err := db.QueryRow("SELECT value FROM public.variables WHERE key=$1;", key).Scan(&s)
	if err != nil {
		return "", err
	}

	return s, nil
}

// RemoveVar deletes an entry in the variables table.
func (db *DBM) RemoveVar(key string) error {
	_, err := db.Exec("DELETE FROM public.variables WHERE key=$1;", key)
	return err
}
