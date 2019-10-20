package anthropoi

import (
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User account structure holds basic login and personal information.
type User struct {
	/*
	 * Required bits
	 */

	// ID of user in the database.
	ID int64
	// Username to log in with.
	Usermame string
	// Password for user account.
	Password string
	// Salt for the password.
	Salt string
	// Email to verify account or reset password.
	Email string
	// Created timestamp.
	Created time.Time
	// Locked accounts can't log in.
	Locked bool

	/*
	 * Optional bits
	 */

	// Sites the user is a member of.
	Sites []string

	// First name of user (optional).
	First string
	// Last name of user (optional).
	Last string
	// Data for the account. JSON field for all the customising you need.
	Data string
	// Tokens is meant to store any authentication tokens required for external sites.
	Tokens string
}

// AddUser creates a new User. This may fail.
func (db *DBM) AddUser(username, password, email, first, last, data, tokens string, cost int) (*User, error) {
	u := &User{
		Usermame: username,
		Email:    email,
		First:    first,
		Last:     last,
		Data:     data,
		Tokens:   tokens,
	}

	if cost < 10 {
		cost = 10
	}

	err := u.SetPassword(password, cost)
	if err != nil {
		return nil, err
	}

	q := "INSERT INTO public.users (username,password,salt,email,first,last,data,tokens) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id;"
	st, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}

	defer st.Close()
	err = st.QueryRow(u.Usermame, u.Password, u.Salt, u.Email, u.First, u.Last, u.Data, u.Tokens).Scan(&u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateUser saves an existing user by ID.
func (db *DBM) SaveUser(u *User) error {
	if u.Data == "" {
		u.Data = "{}"
	}

	if u.Tokens == "" {
		u.Tokens = "{}"
	}

	q := `UPDATE public.users SET username=$1,password=$2,salt=$3,email=$4,locked=$5,first=$6,last=$7,data=$8,tokens=$9 WHERE id=$10;`
	_, err := db.Exec(q, u.Usermame, u.Password, u.Salt, u.Email, u.Locked, u.First, u.Last, u.Data, u.Tokens, u.ID)
	if err != nil {
		fmt.Printf("WTF? %s\n", err.Error())
		return err
	}

	return nil
}

// GetUser returns a User based on an ID.
func (db *DBM) GetUser(id int64) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens FROM public.users WHERE id=$1 LIMIT 1", id).Scan(
		&u.ID, &u.Usermame, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens)
	if err != nil {
		return nil, err
	}

	return &u, db.GetSitesForUser(&u)
}

// GetUserByName for when you don't have an ID.
func (db *DBM) GetUserByName(name string) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id,username,password,salt,email,created,locked,first,last,data,tokens FROM public.users WHERE username=$1 LIMIT 1", name).Scan(
		&u.ID, &u.Usermame, &u.Password, &u.Salt, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last, &u.Data, &u.Tokens)
	if err != nil {
		return nil, err
	}

	return &u, db.GetSitesForUser(&u)
}

// GetSitesForUser fills the Sites field in the User struct.
func (db *DBM) GetSitesForUser(u *User) error {
	q := `SELECT name FROM public.users u
	INNER JOIN membership m ON u.id=m.userid
	INNER JOIN sites s ON m.siteid=s.id WHERE u.id=$1;`
	rows, err := db.Query(q, u.ID)
	if err != nil {
		return err
	}

	defer rows.Close()
	u.Sites = []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		u.Sites = append(u.Sites, name)
	}
	return nil
}

// DeleteUser by ID.
func (db *DBM) DeleteUser(id int64) error {
	_, err := db.Exec("DELETE FROM public.users WHERE id=$1;", id)
	return err
}

// DeleteUserByName for when that's needed.
func (db *DBM) DeleteUserByName(name string) error {
	_, err := db.Exec("DELETE FROM public.users WHERE username=$1;", name)
	return err
}

// GetUsers retrieves all users, up to a limit, sorted by ID.
func (db *DBM) GetUsers(limit int64) ([]*User, error) {
	q := "SELECT id,username,email,created,locked,first,last FROM public.users"
	if limit > 0 {
		q += " LIMIT $1"
	}

	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = db.Query(q, limit)
	} else {
		rows, err = db.Query(q)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var list []*User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Usermame, &u.Email, &u.Created, &u.Locked, &u.First, &u.Last)
		if err != nil {
			return nil, err
		}

		err = db.GetSitesForUser(&u)
		if err != nil {
			return nil, err
		}

		list = append(list, &u)
	}
	return list, nil
}

// SetPassword generates a new salt and sets the password.
func (u *User) SetPassword(password string, cost int) error {
	if cost < 10 {
		cost = 10
	}

	u.Salt = GenString(32)
	s := password + u.Salt
	hash, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// GenerateDovecotPassword creates a Dovecot-compatible password with the SHA512-CRYPT algorithm prefix.
func GenerateDovecotPassword(password, salt string, rounds int) string {
	if rounds == 0 {
		rounds = 100000
	}

	alt := sha512.New()
	alt.Write([]byte(password))
	alt.Write([]byte(salt))
	alt.Write([]byte(password))
	altsum := alt.Sum(nil)

	intermediate := sha512.New()
	intermediate.Write([]byte(password))
	intermediate.Write([]byte(salt))
	var l int
	for l = len(password); l > 64; l -= 64 {
		intermediate.Write([]byte(altsum))
	}
	intermediate.Write(altsum[:l])

	for l = len(password); l > 0; l >>= 1 {
		if (l & 1) == 0 {
			intermediate.Write([]byte(password))
		} else {
			intermediate.Write(altsum)
		}
	}
	isum := intermediate.Sum(nil)

	// S bytes
	S := sha512.New()
	for i := 0; i < (16 + int(isum[0])); i++ {
		S.Write([]byte(salt))
	}
	Ssum := S.Sum(nil)
	Sseq := make([]byte, 0, len(salt))
	for l = len(salt); l > 64; l -= 64 {
		Sseq = append(Sseq, Ssum...)
	}
	Sseq = append(Sseq, Ssum[:l]...)

	// P bytes
	P := sha512.New()
	for i := 0; i < len(password); i++ {
		P.Write([]byte(password))
	}
	Psum := P.Sum(nil)
	Pseq := make([]byte, 0, len(password))
	for l = len(password); l > 64; l -= 64 {
		Pseq = append(Pseq, Psum...)
	}
	Pseq = append(Pseq, Psum[:l]...)

	sum := isum
	for i := 0; i < rounds; i++ {
		hash := sha512.New()
		if (i & 1) != 0 {

			hash.Write(Pseq)
		} else {
			hash.Write(sum)
		}

		if (i % 3) != 0 {
			hash.Write(Sseq)
		}

		if (i % 7) != 0 {
			hash.Write(Pseq)
		}

		if (i & 1) != 0 {
			hash.Write(sum)
		} else {
			hash.Write(Pseq)
		}

		sum = hash.Sum(nil)
	}

	in := []byte{sum[42], sum[21], sum[0],
		sum[1], sum[43], sum[22],
		sum[23], sum[2], sum[44],
		sum[45], sum[24], sum[3],
		sum[4], sum[46], sum[25],
		sum[26], sum[5], sum[47],
		sum[48], sum[27], sum[6],
		sum[7], sum[49], sum[28],
		sum[29], sum[8], sum[50],
		sum[51], sum[30], sum[9],
		sum[10], sum[52], sum[31],
		sum[32], sum[11], sum[53],
		sum[54], sum[33], sum[12],
		sum[13], sum[55], sum[34],
		sum[35], sum[14], sum[56],
		sum[57], sum[36], sum[15],
		sum[16], sum[58], sum[37],
		sum[38], sum[17], sum[59],
		sum[60], sum[39], sum[18],
		sum[19], sum[61], sum[40],
		sum[41], sum[20], sum[62],
		sum[63]}
	return fmt.Sprintf("{SHA512-CRYPT}$6$rounds=%d$%s$%s", rounds, salt, Base6424(string(in)))
}

// SetDovecotPassword sets a Dovecot IMAP-compatible password for the user.
func (u *User) SetDovecotPassword(password string, rounds int) {
	u.Salt = GenString(16)
	u.Password = GenerateDovecotPassword(password, u.Salt, rounds)
}

// CheckPassword against the account's hash.
func (u *User) CheckPassword(password string) bool {
	if u.Usermame == "" || u.Password == "" {
		return false
	}

	if u.Password[:4] == pre_bcrypt {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		return (err == nil)
	}

	if strings.HasPrefix(u.Password, pre_sha512crypt) {

	}
	return false
}

func (u *User) CompareDovecotHashAndPassword(password string) bool {
	a := strings.Split(u.Password, "$")
	if len(a) != 5 {
		return false
	}

	ra := strings.Split(a[2], "=")
	if len(ra) != 2 {
		return false
	}

	rounds, err := strconv.Atoi(ra[1])
	if err != nil {
		return false
	}
	pw := GenerateDovecotPassword(password, u.Salt, rounds)
	a2 := strings.Split(pw, "$")
	if len(a2) != 5 {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(a[4]), []byte(a2[4])) == 1
}
