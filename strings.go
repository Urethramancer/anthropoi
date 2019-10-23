package anthropoi

import (
	"crypto/rand"
	"math/big"
)

const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'#$%&/()=?@*^<>-.:,;|[]{}"

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
