package anthropoi

// SetFlag sets a string in the flags table.
func (db *DBM) SetFlag(key string, flag bool) error {
	q := `INSERT INTO public.flags (key,flag) 
	VALUES ($1, $2)
	ON CONFLICT (key) DO UPDATE SET flag = $2;`
	_, err := db.Exec(q, key, flag)
	return err
}

// GetFlag gets a string from the flags table.
func (db *DBM) GetFlag(key string) bool {
	var f bool
	err := db.QueryRow("SELECT flag FROM public.flags WHERE key=$1;", key).Scan(&f)
	if err != nil {
		// No key is the same as false.
		return false
	}

	return f
}

// ClearFlag removes an entry in the flags table.
func (db *DBM) ClearFlag(key string) error {
	_, err := db.Exec("DELETE FROM public.flags WHERE key=$1;", key)
	return err
}
