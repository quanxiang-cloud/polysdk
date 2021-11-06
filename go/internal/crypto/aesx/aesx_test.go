package aesx

import (
	"fmt"
	"testing"
)

func TestAesx(t *testing.T) {
	txt := "foo"
	keys := []string{"foo", "bar"}
	enc, _ := EncodeString(txt, keys...)
	fmt.Println(enc)
	dec, _ := DecodeString(enc, keys...)
	fmt.Println(dec)
}
