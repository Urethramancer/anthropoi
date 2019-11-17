package anthropoi

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// AddResetToken creates a reset token with an expiry for an account, then returns the human-readable hash.
func (db *DBM) AddResetToken(user *User, duration time.Duration) (string, error) {
	future := time.Now().Add(duration)
	h := sha256.New()
	h.Write([]byte(user.Username))
	h.Write([]byte(GenString(16)))
	h.Write([]byte(user.Salt))
	hash := hex.EncodeToString(h.Sum(nil))

	q := "INSERT INTO public.resetkeys (key,account,expiry) VALUES($1,$2,$3);"
	_, err := db.Exec(q, hash, user.ID, future)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// DeleteResetToken invalidates a reset token. Call when resetting a password.
func (db *DBM) DeleteResetToken(key string) error {
	_, err := db.Exec("DELETE FROM public.resetkeys WHERE key=$1", key)
	return err
}

// GetUserForReset returns a User if the token is valid.
func (db *DBM) GetUserForReset(key string) (*User, error) {
	q := `SELECT id,username,password,salt,email,created,locked,first,last,data,
	tokens,admin FROM public.users u
	INNER JOIN public.resetkeys r ON u.id=r.account
	WHERE r.key=%1;`
	var user User
	err := db.QueryRow(q, key).Scan(&user)
	return &user, err
}

// PurgeResetTokens deletes the oldest unused and expired tokens.
func (db *DBM) PurgeResetTokens() error {
	_, err := db.Exec("DELETE FROM public.resetkeys WHERE expiry<now();")
	return err
}
