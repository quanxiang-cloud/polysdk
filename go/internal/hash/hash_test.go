package hash

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	b := bytes.NewBuffer(make([]byte, 0, 8))
	if err := binary.Write(b, binary.LittleEndian, uint64(256+5)); err != nil {
		panic(err)
	}
	if bb := b.Bytes(); !reflect.DeepEqual(bb, []byte{5, 1, 0, 0, 0, 0, 0, 0}) {
		t.Fatal(bb)
	}
	if !reflect.DeepEqual(hashSep, []byte{0, 59, 13, 9, 10}) {
		panic(hashSep)
	}

	//--------------------------------------------------------------------------

	type testCase struct {
		name   string
		src    []string
		expect string
	}
	cases := []*testCase{
		&testCase{
			src:    []string{"foo", "bar"},
			expect: "7f4eeee9b6bda8d679516686b8f4ef4721fd2515e8b204ed832317e733151454",
		},
		&testCase{
			src:    []string{"foobar"},
			expect: "40ab974761e4d9ab9e1067669b137224dd2a0df85ce525b6ea0288be55b3e6fe",
		},
	}

	for i, v := range cases {
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d TestHash(%s) fail: expect %v, got %v", i+1, msg, expect, got)
			}
		}
		got := DefaultHexString(v.src...)
		assert(len(got), DefaultSize()*2, v.name)
		assert(got, v.expect, v.name)
	}
}
