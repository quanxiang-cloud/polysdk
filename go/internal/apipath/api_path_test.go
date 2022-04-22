package apipath

import (
	"fmt"
	"reflect"
	"testing"
)

const testPrint = false

func TestJoin(t *testing.T) {
	type testCase struct {
		ns     string
		name   string
		expect string
	}
	cases := []*testCase{
		&testCase{ns: "", name: "a", expect: "/a"},
		&testCase{ns: "/a", name: "b", expect: "/a/b"},
		&testCase{ns: "a", name: "b", expect: "/a/b"},
	}
	for i, v := range cases {
		got := Join(v.ns, v.name)
		if testPrint && true {
			fmt.Printf("Join(\"%s\", \"%s\")=\"%s\"\n", v.ns, v.name, got)
		}
		if got != v.expect {
			t.Errorf("case %d ns=%s name=%s expect=%s got=%s",
				i+1, v.ns, v.name, v.expect, got)
		}
	}
}

func TestSplit(t *testing.T) {
	type testCase struct {
		full string
		ns   string
		name string
	}
	cases := []*testCase{
		&testCase{full: "a", ns: "", name: "a"},
		&testCase{full: "/a", ns: "", name: "a"},
		&testCase{full: "a/b", ns: "/a", name: "b"},
		&testCase{full: "/a/b", ns: "/a", name: "b"},
	}
	for i, v := range cases {
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d TestSplit(%s) fail: expect %v, got %v", i+1, msg, expect, got)
			}
		}
		ns, name := Split(v.full)
		if testPrint && true {
			fmt.Printf("Split(\"%s\")=\"%s\", \"%s\"\n", v.full, v.ns, v.name)
		}
		assert(Name(v.full), name, "Name()")
		assert(Parent(v.full), ns, "Parent()")
		if ns != v.ns || name != v.name {
			t.Errorf("case %d full=%s ns=%s/%s name=%s/%s, mismatch",
				i+1, v.full, ns, v.ns, name, v.name)
		}
	}
}

func TestFormat(t *testing.T) {
	type testCase struct {
		full   string
		expect string
	}
	cases := []*testCase{
		&testCase{full: "", expect: "/"},
		&testCase{full: "a", expect: "/a"},
		&testCase{full: "/a", expect: "/a"},
		&testCase{full: "a/b", expect: "/a/b"},
		&testCase{full: "/a/b", expect: "/a/b"},
	}
	for i, v := range cases {
		got := Format(v.full)
		if testPrint && true {
			fmt.Printf("Format(\"%s\")=\"%s\", \"%s\"\n", v.full, v.expect, got)
		}
		if got != v.expect {
			t.Errorf("case %d full=%s format=%s/%s, mismatch",
				i+1, v.full, got, v.expect)
		}
	}
}
