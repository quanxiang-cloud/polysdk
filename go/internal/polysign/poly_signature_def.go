// Package polysign is signature defines for polyapi
package polysign

import (
	"encoding/json"
)

// signature header
const (
	XHeaderPolySignVersion   = "X-Polysign-Version"
	XHeaderPolySignMethod    = "X-Polysign-Method"
	XHeaderPolySignKeyID     = "X-Polysign-Access-Key-Id"
	XHeaderPolySignTimestamp = "X-Polysign-Timestamp"
	XHeaderPolySignSignature = "X-Polysign-Signature" // NOTE: client signature result
)

// signature header value
const (
	XHeaderPolySignVersionVal   = "1"
	XHeaderPolySignMethodVal    = "HmacSHA256"
	ISO8601                     = "2006-01-02T15:04:05-0700" // ISO8601
	XHeaderPolySignTimestampFmt = ISO8601
)

// special body field define
const (
	// XPolyBodyHideArgs is poly reserve field in body
	// NOTE: pass path arg of raw api by this object
	XPolyBodyHideArgs = "$polyapi_hide$"

	// NOTE: this name means this is real body root of customer api
	XPolyCustomerBodyRoot = "$body$"

	// XPolyRaiseUpFieldName is a special filed name.
	// NOTE: if a field with this name, generate query will raiseup its children
	// eg: {"a":1,"b":2} is the same as {"a":1,"$$*_*$$":{"b":2}}
	XPolyRaiseUpFieldName = "$$*_*$$"
)

// PolySignatureInfo is the data structure for signature generator
type PolySignatureInfo struct {
	AccessKeyID string `json:"X-Polysign-Access-Key-Id"` // header
	Timestamp   string `json:"X-Polysign-Timestamp"`     // header
	SignMethod  string `json:"X-Polysign-Method"`        // header
	SignVersion string `json:"X-Polysign-Version"`       // header

	// NOTE: body XPolyIgnoreFieldName defined, signature will ignore name for this field
	// Body:{Child:foo} will generate as Child=foo in query
	Body json.RawMessage `json:"$$*_*$$"`
}
