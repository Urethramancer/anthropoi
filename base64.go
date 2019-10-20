package anthropoi

import (
	"strings"
)

const alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base6424 used by some password hashing algorithms.
func Base6424(src string) string {
	if len(src) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for len(src) > 0 {
		switch len(src) {
		default:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[((src[0]>>6)|(src[1]<<2))&0x3f])
			_ = buf.WriteByte(alphabet[((src[1]>>4)|(src[2]<<4))&0x3f])
			_ = buf.WriteByte(alphabet[(src[2]>>2)&0x3f])
			src = src[3:]
		case 2:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[((src[0]>>6)|(src[1]<<2))&0x3f])
			_ = buf.WriteByte(alphabet[(src[1]>>4)&0x3f])
			src = src[2:]
		case 1:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[(src[0]>>6)&0x3f])
			src = src[1:]
		}
	}
	return buf.String()
}
