package polysign

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRaiseUpFiledTag(t *testing.T) {
	body := `this is raise up field(body)`
	var testCase = PolySignatureInfo{
		Signature:   "this is XBodyPolySignSignature",
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
	verify(XBodyPolySignSignature, testCase.Signature)
	verify(XPolyRaiseUpFieldName, body)
}
