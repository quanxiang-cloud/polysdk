package aesx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"polysdk/consts"
	"polysdk/internal/hash"
	"time"
)

var (
	saltP          = "#jK‥y8&%d\\nf%$ǚ$*)dK\u000038^%3☒3 \t?2i\n/L-<9U"
	saltS          = "UyH-j*<&^-o‥p*6ǘ5%/HJ^k,k\u000086\t\n4j ǎguI$#jg\\nhj@!jkJhl"
	base64Encoding = base64.RawURLEncoding
	random         = rand.Read
)

const (
	ivSize            = aes.BlockSize
	appendSize        = 1
	version      byte = 1
	cryptoRounds      = 30_000_000 * consts.HashTimeEstimateMS / int(time.Second/time.Millisecond) //2^25
)

// EncodeString encode a string with strong aes crypto.
// NOTE: It is deliberately to make slow to avoid attack.
func EncodeString(plantText string, keys ...string) (string, error) {
	if plantText == "" {
		return "", errors.New("empty plant text")
	}
	if len(keys) == 0 {
		return "", errors.New("empty keys")
	}
	b := []byte(plantText)
	full, iv, data := makeBuffer(len(b))
	_, err := random(iv)
	if err != nil {
		return "", err
	}
	c, err := newCipher(keys, iv)
	if err != nil {
		return "", err
	}
	data[0] = version
	if n := copy(data[appendSize:], b); n != len(b) {
		return "", errors.New("encrypt fail")
	}
	for i := 0; i < cryptoRounds; i++ {
		c.XORKeyStream(data, data)
	}
	return base64Encoding.EncodeToString(full), nil
}

// DecodeString decode a string with strong aes crypto.
// NOTE: It is deliberately to make slow to avoid attack.
func DecodeString(encryptText string, keys ...string) (string, error) {
	if encryptText == "" {
		return "", errors.New("empty crypto text")
	}
	if len(keys) == 0 {
		return "", errors.New("empty keys")
	}

	iv, data, err := toBuffer(encryptText)
	if err != nil {
		return "", err
	}

	c, err := newCipher(keys, iv)
	if err != nil {
		return "", err
	}

	for i := 0; i < cryptoRounds; i++ {
		c.XORKeyStream(data, data)
	}
	if data[0] != version {
		return "", errors.New("decode fail")
	}
	return string(data[appendSize:]), nil
}

//------------------------------------------------------------------------------

func makeBuffer(dataSize int) ([]byte, []byte, []byte) {
	if dataSize < 0 {
		panic(dataSize)
	}
	b := make([]byte, ivSize+dataSize+appendSize)
	full, iv, data := b, b[:ivSize], b[ivSize:]

	return full, iv, data
}

func toBuffer(encryptText string) ([]byte, []byte, error) {
	b, err := base64Encoding.DecodeString(encryptText)
	if err != nil {
		return nil, nil, err
	}
	if len(b) <= ivSize+appendSize {
		return nil, nil, errors.New("encrypt text too short")
	}
	iv, data := b[:ivSize], b[ivSize:]
	return iv, data, nil
}

func newCipher(keys []string, iv []byte) (cipher.Stream, error) {
	aesKey := hashKey(nil, keys...)
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	cc := cipher.NewCTR(c, iv)
	return cc, nil
}

func hashKey(buf []byte, elems ...string) []byte {
	keys := make([]string, 0, len(elems)+2)
	keys = append(keys, saltP)
	keys = append(keys, elems...)
	keys = append(keys, saltS)

	b := hash.Default(buf, keys...)
	b[(b[0]>>3)&0x0F] ^= 0x5A // non-standard hash
	return b
}
