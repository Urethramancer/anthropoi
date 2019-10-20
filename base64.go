package anthropoi

const alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base6424 used by some password hashing algorithms.
func Base6424(src []byte) (hash []byte) {
	if len(src) == 0 {
		return []byte{} // TODO: return nil
	}

	hashSize := (len(src) * 8) / 6
	if (len(src) % 6) != 0 {
		hashSize += 1
	}
	hash = make([]byte, hashSize)

	dst := hash
	for len(src) > 0 {
		switch len(src) {
		default:
			dst[0] = alphabet[src[0]&0x3f]
			dst[1] = alphabet[((src[0]>>6)|(src[1]<<2))&0x3f]
			dst[2] = alphabet[((src[1]>>4)|(src[2]<<4))&0x3f]
			dst[3] = alphabet[(src[2]>>2)&0x3f]
			src = src[3:]
			dst = dst[4:]
		case 2:
			dst[0] = alphabet[src[0]&0x3f]
			dst[1] = alphabet[((src[0]>>6)|(src[1]<<2))&0x3f]
			dst[2] = alphabet[(src[1]>>4)&0x3f]
			src = src[2:]
			dst = dst[3:]
		case 1:
			dst[0] = alphabet[src[0]&0x3f]
			dst[1] = alphabet[(src[0]>>6)&0x3f]
			src = src[1:]
			dst = dst[2:]
		}
	}

	return
}
