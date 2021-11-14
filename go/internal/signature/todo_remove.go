// TODO: remove this file*************

package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Sha256 Sha256
func Sha256(entity []byte, key []byte) ([]byte, error) {
	return sha(entity, hmac.New(sha256.New, key))
}

// Signature generator
func Signature(src interface{}, key string) (string, error) {
	query := dismantling("", src)

	body, err := Sha256([]byte(query), []byte(key))
	if err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(body)

	return signature, nil
}

func _toRawQuery(src interface{}) string {
	b, ok := src.([]byte)
	if !ok {
		return ""
	}
	var d interface{}
	if err := json.Unmarshal(b, &d); err != nil {
		return err.Error()
	}
	//fmt.Printf("%#v\n", d)
	return dismantling("", d)
}

func dismantlingBase(key string, value interface{}) string {
	buf := bytes.Buffer{}
	buf.WriteString(key)
	buf.WriteString("=")
	val := fmt.Sprint(value)
	val = url.QueryEscape(val)
	buf.WriteString(val)
	return buf.String()
}
func dismantlingMap(key string, value interface{}) string {
	buf := bytes.Buffer{}
	key = division(key)

	var keys []string
	valueOfValue := reflect.ValueOf(value)
	numField := valueOfValue.Len()
	keys = make([]string, 0, numField)
	vals := make(map[string]string, numField)

	iter := valueOfValue.MapRange()
	for iter.Next() {
		if !iter.Value().CanInterface() {
			continue
		}
		k := iter.Key().String()
		keys = append(keys, k)

		val := iter.Value().Interface()
		vals[k] = dismantling(key+k, val)
	}

	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString("&")
		buf.WriteString(vals[k])
	}

	return trim(buf.String())
}

func dismantlingArray(key string, value interface{}) string {
	buf := bytes.Buffer{}
	key = division(key)

	valueOfValue := reflect.ValueOf(value)
	numField := valueOfValue.Len()

	for i := 0; i < numField; i++ {
		item := valueOfValue.Index(i)
		if !item.CanInterface() {
			continue
		}
		buf.WriteString("&")
		buf.WriteString(dismantling(key+strconv.Itoa(i+1), item))
	}
	return trim(buf.String())
}

func dismantling(key string, value interface{}) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		return dismantlingMap(key, value)
	case reflect.Array, reflect.Slice:
		return dismantlingArray(key, value)
	default:
	}

	return dismantlingBase(key, value)
}

func division(key string) string {
	if key != "" {
		key = key + "."
	}
	return key
}

func trim(query string) string {
	if strings.HasPrefix(query, "&") {
		return query[1:]
	}
	return query
}
