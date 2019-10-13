package anthropoi

import (
	"crypto/rand"
	"math/big"
)

const (
	DefaultName         = "accounts"
	databaseDefinitions = `CREATE DATABASE {NAME};`
	databaseTriggers    = `BEGIN WORK;
	-- We'll trigger creation timestamp setting in a few places.
	CREATE OR REPLACE FUNCTION trigger_set_timestamp()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.created = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	COMMIT WORK;`
	validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'#$%&/()=?@*^<>-.:,;|[]{}"
)

// GenString generates a random string, usable for passwords.
func GenString(size int) string {
	s := make([]byte, size)
	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(validChars))))
		if err != nil {
			return ""
		}
		c := validChars[n.Int64()]
		s[i] = c
	}
	return string(s)
}
