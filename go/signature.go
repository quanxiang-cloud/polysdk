package polysdk

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// const parameter for body._signature
const (
	signatureVersion = 1
	signatureMethod  = "HmacSHA256"

	TimeISO8601 = "2006-01-02T15:04:05-0700" // ISO8601 timestamp format
)

type bodySignature struct {
	Version     uint64 `json:"version"`
	Method      string `json:"method"`
	AccessKeyID string `json:"access_key_id"`
	Timestamp   string `json:"timestamp"`
}

func newSignatureTemplate(accessKeyID string) (*signatureTemplate, error) {
	tmp := bodySignature{
		Version:     signatureVersion,
		Method:      signatureMethod,
		AccessKeyID: accessKeyID,
		Timestamp:   TimeISO8601,
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	st := string(b)
	t := &signatureTemplate{
		signatureTemplate: st,
		timestampIndex:    strings.Index(st, TimeISO8601),
	}
	if t.timestampIndex < 0 {
		return nil, errors.New("failed to find timestamp in template")
	}

	return t, nil
}

// signatureTemplate use template to replace by timestamp()
type signatureTemplate struct {
	// {"version":1,"method":"HmacSHA256","access_key_id":"$access_key_id$","timestamp":"2006-01-02T15:04:05-0700"}
	signatureTemplate string
	// 82
	timestampIndex int
}

// replace timestamp in template and return new bytes
func (t *signatureTemplate) genBodySignature() json.RawMessage {
	b := []byte(t.signatureTemplate)
	toReplace := b[t.timestampIndex : t.timestampIndex+len(TimeISO8601)]
	copy(toReplace, []byte(timestamp())) // replace timestamp in template
	return json.RawMessage(b)
}

func timestamp() string {
	return time.Now().Format(TimeISO8601)
}

// Timestamp get current timestamp of UTC
func Timestamp(f string) string {
	if f == "" {
		f = "2006-01-02T15:04:05Z"
	}
	return time.Now().UTC().Format(f)
}
