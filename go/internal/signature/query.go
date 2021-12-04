package signature

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"polysdk/internal/polysign"
	"sort"
)

// ToQuery convert any date to http GET query parameter with ordered.
func ToQuery(data interface{}) string {
	buf := bytes.NewBuffer(nil)
	if err := buildQuery("", data, buf, 0); err != nil {
		return ""
	}
	return trimQuery(buf.String())
}

func buildQuery(name string, d interface{}, buf *bytes.Buffer, depth int) error {
	if depth >= 20 {
		return errors.New("buildQuery out of recursion")
	}
	writeSingle := func(v interface{}) {
		buf.WriteString(fmt.Sprintf("&%s=", name))
		buf.WriteString(url.QueryEscape(fmt.Sprint(v)))
	}
	fromBytes := func(v []byte) error {
		var x interface{}
		if err := json.Unmarshal(v, &x); err != nil {
			return err
		}
		return buildQuery(name, x, buf, depth) // NOTE: dont increase depth here
	}
	switch v := d.(type) {
	case string:
		writeSingle(v)
	case float64:
		writeSingle(v)
	case bool:
		writeSingle(v)
	case map[string]interface{}:
		if depth == 0 {
			// NOTE: dont treat raiseField as a child, raise up it's children to parent
			const raiseField = polysign.XPolyRaiseUpFieldName
			if raise, ok := v[raiseField]; ok {
				delete(v, raiseField) //remove from parent
				if mp, ok := raise.(map[string]interface{}); ok {
					for k, c := range mp {
						// TODO: handle duplication?
						// if _, ok := v[k]; ok {
						// 	return fmt.Errorf("duplicate field %s between root and raise up fields", k)
						// }
						v[k] = c
					}
				} else {
					v[polysign.XPolyCustomerBodyRoot] = raise
				}
			}
		}

		names := make([]string, 0, len(v))
		for k := range v {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, k := range names {
			n := k
			if name != "" {
				n = fmt.Sprintf("%s.%s", name, k)
			}
			if val := v[k]; val != nil {
				if err := buildQuery(n, v[k], buf, depth+1); err != nil {
					return err
				}
			}
		}
	case []interface{}:
		if name != "" {
			for i, vv := range v {
				n := fmt.Sprintf("%s.%d", name, i+1)
				buildQuery(n, vv, buf, depth+1)
			}
		}
	//--------------------------------------------------------------------------
	default:
		b, err := json.Marshal(d)
		if err != nil {
			return err
		}
		return fromBytes(b)
	case []byte:
		return fromBytes(v)
	case json.RawMessage:
		return fromBytes([]byte(v))
	}
	return nil
}

func trimQuery(query string) string {
	if len(query) > 0 && query[0] == '&' {
		return query[1:]
	}
	return query
}
