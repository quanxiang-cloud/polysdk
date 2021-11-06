package hash

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	hashRounds   = 524288 / 1000 // defatut hash rounds, 2^19
	hashIndex    = 97351831      // default salt index
	hashSaltName = "poly_salt_fa391be3"
)

var (
	hashSep = []byte("\u0000;\r\t\n") //seprater between keys [0 59 13 9 10]
)

// HexString return hex encoding string of []byte
func HexString(b []byte) string {
	return hex.EncodeToString(b)
}

// DefaultHexString generate default hex encoded hash.
// NOTE: It is deliberately to make slow to avoid hash conflict attack.
func DefaultHexString(elems ...string) string {
	h := Default(nil, elems...)
	return HexString(h)
}

// Default generate default hash bytes.
// NOTE: It is deliberately to make slow to avoid hash conflict attack.
func Default(buf []byte, elems ...string) []byte {
	h := Sha256Hash(buf, hashRounds, hashIndex, elems...)
	return h
}

// DefaultSize return size of defalut hash
func DefaultSize() int {
	return sha256.Size
}

// Sha256Hash generate sha256 hash.
// if index>0, it try to avoid hash conflict by salt
// round is used to control difficulty of hash conflict
// NOTE: It panic when round <= 0 or missing elems
func Sha256Hash(buf []byte, round int, index int, elems ...string) []byte {
	if len(elems) <= 0 {
		panic("missing hash key")
	}
	if round <= 0 {
		panic("round must >0")
	}

	h := sha256.New()

	src := make([][]byte, 0, len(elems)+1)
	for _, v := range elems {
		src = append(src, []byte(v))
	}
	if index > 0 {
		src = append(src, []byte(fmt.Sprintf("%s_%d", hashSaltName, index)))
	}

	if len(buf) < h.Size() {
		buf = make([]byte, h.Size())
	}
	for i := 0; i == 0 || i < round; i++ { // at least loop once
		// write round
		if err := binary.Write(h, binary.LittleEndian, uint64(i)); err != nil {
			panic(err)
		}

		for j, v := range src { // write keys
			// write key index
			if err := binary.Write(h, binary.LittleEndian, uint64(j)); err != nil {
				panic(err)
			}

			if _, err := h.Write(v); err != nil { // write key
				panic(err)
			}
			if _, err := h.Write(hashSep); err != nil { // write key seprater
				panic(err)
			}
		}
		if _, err := h.Write(h.Sum(buf[:0])); err != nil { // write the hash sum self
			panic(err)
		}
	}

	return h.Sum(buf[:0])
}
