package polysign

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIgnoreFiledTag(t *testing.T) {
	body := `this is ignore name field(body)`
	var testCase = PolySignatureInfo{
		AccessKeyID: "this is XHeaderPolySignKeyID",
		Timestamp:   "this is XHeaderPolySignTimestamp",
		SignMethod:  "this is XHeaderPolySignMethod",
		SignVersion: "this is XHeaderPolySignVersion",
		Body:        []byte(fmt.Sprintf(`"%s"`, body)),
	}
	b, err := json.MarshalIndent(testCase, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	var got map[string]string
	if err := json.Unmarshal(b, &got); err != nil {
		panic(err)
	}
	verify := func(name string, expect string) {
		s, ok := got[name]
		if !ok {
			t.Errorf("missing filed %s", name)
		}
		if s != expect {
			t.Errorf("field %s mismatch. expect=%q got=%q", name, expect, s)
		}
	}
	verify(XHeaderPolySignVersion, testCase.SignVersion)
	verify(XHeaderPolySignMethod, testCase.SignMethod)
	verify(XHeaderPolySignKeyID, testCase.AccessKeyID)
	verify(XHeaderPolySignTimestamp, testCase.Timestamp)
	verify(XPolyRaiseUpFieldName, body)
}
