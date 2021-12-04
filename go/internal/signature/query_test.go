package signature

import (
	"encoding/json"
	"fmt"
	"polysdk/internal/config"
	"polysdk/internal/polysign"
	"testing"
)

func TestBodyToQuery(t *testing.T) {
	var v = map[string]interface{}{
		"a": "foo",
		"b": []string{"foo", "bar"},
		"c": 123,
		"d": 123.456,
		"e": true,
		"f": map[string]interface{}{
			"x":  "xx",
			"y":  []string{"foo2", "bar2"},
			"y2": []float64{3, 4},
			"y3": []bool{true, false},
			"z":  123.4,
			"x2": true,
		},
	}
	b, _ := json.Marshal(v)
	s := string(b)
	fmt.Println(s)
	q := ToQuery(v)
	fmt.Println(q)
	fmt.Println(dismantling("", v))
	var st = struct {
		X  string   `json:"x"`
		Y  []string `json:"y"`
		Y2 []int    `json:"y2"`
		Y3 []bool   `json:"y3"`
		Z  float32  `json:"z"`
		X2 bool     `json:"x2"`
	}{
		X:  "xx",
		Y:  []string{"foo2", "bar2"},
		Y2: []int{3, 4},
		Y3: []bool{true, false},
		Z:  123.4,
		X2: true,
	}
	fmt.Println(ToQuery(v["f"]))
	fmt.Println(dismantling("", v["f"]))
	fmt.Println(ToQuery(st))
	fmt.Println(dismantling("", st))
}

func TestSignature(t *testing.T) {
	demo := map[string]interface{}{
		"id":     "1",
		"name":   "张三",
		"age":    18,
		"status": true,
		"address": map[string]interface{}{
			"contry":   "china",
			"province": "sichuang",
			"city":     "chengdu",
		},
		"interest":  []string{"basketball", "football", "pingpong"},
		"timestamp": "2021-01-29T17:23:05-0800",
	}

	secretKey := "foo"
	cryptedKey := config.Encrypt(secretKey)
	println("cryptedKey", cryptedKey)
	s, err := NewSigner(cryptedKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Signature(demo))
	fmt.Println(s.Signature(demo))
	fmt.Println(Signature(demo, secretKey))
	jsonBytes := []byte(`{"x":"foo","y":"bar"}`)
	fmt.Println(s.Signature(jsonBytes))
	x := map[string]interface{}{
		"x": "foo",
		"y": "bar",
	}
	fmt.Println(s.Signature(x))
	fmt.Println(Signature(x, secretKey))
	fmt.Println(Signature(jsonBytes, secretKey))
}

func TestToQuery(t *testing.T) {
	var testCase = map[string]interface{}{
		"a": "f",
		"b": []string{"foo", "bar"},
		"c": map[string]interface{}{
			"c1": 1,
			"c2": map[string]interface{}{
				"c21": "foo",
				"c22": []interface{}{
					[]interface{}{1, "foo"},
					map[string]interface{}{
						"c221": "foo",
						"c222": true,
					},
				},
			},
		},
		"d": []interface{}{
			map[string]interface{}{
				"c21": "foo",
				"c22": []interface{}{
					[]interface{}{1, "foo"},
					map[string]interface{}{
						"c221": "foo",
						"c222": true,
					},
				},
			},
			[]interface{}{
				[]interface{}{1, "foo"},
				map[string]interface{}{
					"c221": []string{"foo", "bar"},
					"c222": true,
				},
			},
		},
	}

	b, err := json.MarshalIndent(testCase, "", "  ")
	if err != nil {
		panic(err)
	}

	var d interface{}
	if err := json.Unmarshal(b, &d); err != nil {
		panic(err)
	}

	query := ToQuery(d)
	expect := `a=f&b.1=foo&b.2=bar&c.c1=1&c.c2.c21=foo&c.c2.c22.1.1=1&c.c2.c22.1.2=foo&c.c2.c22.2.c221=foo&c.c2.c22.2.c222=true&d.1.c21=foo&d.1.c22.1.1=1&d.1.c22.1.2=foo&d.1.c22.2.c221=foo&d.1.c22.2.c222=true&d.2.1.1=1&d.2.1.2=foo&d.2.2.c221.1=foo&d.2.2.c221.2=bar&d.2.2.c222=true`
	if expect != query {
		fmt.Println(string(b))
		t.Errorf("TestToQuery:\nexpect %s\ngot    %s\n", expect, query)
	}
}

func TestToQueryExt(t *testing.T) {
	type testCase struct {
		body   interface{}
		expect string
	}
	testCases := []*testCase{
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID:  "foo",
				polysign.XPolyRaiseUpFieldName: "stringBody",
			},
			expect: "$body$=stringBody&X-Polysign-Access-Key-Id=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID:  "foo",
				polysign.XPolyRaiseUpFieldName: []string{"foo", "bar"},
			},
			expect: "$body$.1=foo&$body$.2=bar&X-Polysign-Access-Key-Id=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID: "foo",
				polysign.XPolyRaiseUpFieldName: map[string]interface{}{
					"a": "foo",
					polysign.XPolyBodyHideArgs: map[string]interface{}{
						"app": "foo",
					},
					polysign.XPolyCustomerBodyRoot: "bar",
				},
			},
			expect: "$body$=bar&$polyapi_hide$.app=foo&X-Polysign-Access-Key-Id=foo&a=foo",
		},
		&testCase{
			body: map[string]interface{}{
				polysign.XHeaderPolySignKeyID: "foo",
				polysign.XPolyRaiseUpFieldName: map[string]interface{}{
					"a":                           "foo",
					polysign.XHeaderPolySignKeyID: "bar",
				},
			},
			expect: "X-Polysign-Access-Key-Id=bar&a=foo",
		},
	}
	for i, v := range testCases {
		got := ToQuery(v.body)
		if got != v.expect {
			t.Errorf("case %d, expect %q\ngot %q", i+1, v.expect, got)
		}
	}
}
