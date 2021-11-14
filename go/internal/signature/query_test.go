package signature

import (
	"encoding/json"
	"fmt"
	"polysdk/internal/config"
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
