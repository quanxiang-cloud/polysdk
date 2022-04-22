package hash

import (
	"crypto/rand"
)

const (
	// DefaultShortNameLen is defult length of short name
	DefaultShortNameLen = 8
)
const (
	alphaTab   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	trimBits   = 1
	tabLen     = len(alphaTab)
	headTabLen = tabLen - 10 // first byte dont allow number character
)

// ShortID  generate a random string with length n
func ShortID(n int) string {
	s, err := ShortIDWithError(n)
	if err != nil {
		panic(err)
	}
	return s
}

// ShortIDWithError  generate a random string with length n
func ShortIDWithError(n int) (string, error) {
	if n <= 0 {
		n = DefaultShortNameLen
	}
	b := make([]byte, n)
	if nr, err := rand.Read(b); err != nil || nr != len(b) {
		return "", err
	}

	mod := headTabLen
	for i, v := range b {
		idx := (int(v>>trimBits) % mod)
		mod = tabLen

		b[i] = alphaTab[idx]
	}
	return string(b), nil
}
