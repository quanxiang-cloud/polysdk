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
			expect: "7e4f39bab0c9e035c733ed4b19c4147e61116b60222a189ac2e2eec57043e485",
		},
		&testCase{
			src:    []string{"foobar"},
			expect: "e8de1a461e3f686718d1c66ae9874c2fcfaa9119ae75738f17de15c1369a150d",
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
