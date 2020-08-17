package anthropoi

import (
	"crypto/sha512"
	"crypto/subtle"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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

// SetDovecotPassword sets a Dovecot-compatible password for the user.
func (u *User) SetDovecotPassword(password string, rounds int) {
	u.Salt = GenString(16)
	u.Password = GenerateDovecotPassword(password, u.Salt, rounds)
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

// CheckPassword against the account's hash.
func (u *User) CheckPassword(password string) bool {
	if u.Username == "" || u.Password == "" {
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

// CompareDovecotHashAndPassword for systems where getting bcrypt support in Dovecot is a pain.
func (u *User) CompareDovecotHashAndPassword(password string) bool {
	a := u.SplitPasswordElements()
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

// SplitPasswordElements splits the stored password hash and returns it if it fits
// any supported pattern (4 elements for bcrypt, 5 for Dovecot).
func (u *User) SplitPasswordElements() []string {
	a := strings.Split(u.Password, "$")
	if len(a) < 4 || len(a) > 5 {
		return nil
	}

	return a
}

// GetCost for bcrypt hashes.
func (u *User) GetCost() int {
	a := u.SplitPasswordElements()
	if a == nil {
		return 12
	}

	c, err := strconv.Atoi(a[2])
	if err != nil {
		return 12
	}

	return c
}

// GetRounds for Dovecot hashes.
func (u *User) GetRounds() int {
	a := u.SplitPasswordElements()
	if a == nil {
		return 50000
	}

	ar := strings.Split(a[2], "=")
	if len(ar) < 2 {
		return 50000
	}

	r, err := strconv.Atoi(ar[1])
	if err != nil {
		return 50000
	}

	return r
}

// AcceptablePassword does some superficial checking of a potential password.
// It will fail the test if it's too short, contains user details or is all numbers.
// Further policies have to be applied outside of this function.
func (u *User) AcceptablePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	_, err := strconv.ParseInt(password, 10, 64)
	if err == nil {
		return false
	}

	password = strings.ToLower(password)
	comp := strings.ToLower(u.Username)
	if strings.Contains(password, comp) {
		return false
	}

	if strings.Contains(comp, password) {
		return false
	}

	if u.First != "" {
		comp = strings.ToLower(u.First)
		if strings.Contains(password, comp) {
			return false
		}

		if strings.Contains(comp, password) {
			return false
		}
	}

	if u.Last != "" {
		comp = strings.ToLower(u.Last)
		if strings.Contains(password, comp) {
			return false
		}

		if strings.Contains(comp, password) {
			return false
		}
	}

	return true
}
