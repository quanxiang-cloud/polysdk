package polysdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"polysdk/internal/polysign"
	"polysdk/internal/signature"
)

// header define
const (
	// HeaderSignature = "Signature"  // "Signature" in header
	// BodySinature    = "_signature" // _signature in body
	// BodyHide        = "_hide"      // _hide in body

	HeaderContentType = "Content-Type"
)

// Header exports
type Header = http.Header

// http method exports
const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

// BodyBase is base struct of body
type BodyBase struct {
	// {"version":1,"method":"HmacSHA256","access_key_id":"$access_key_id$","timestamp":"2006-01-02T15:04:05-0700"}
	//Signature json.RawMessage `json:"_signature"`

	// path and other hide parameter
	PolyHide map[string]interface{} `json:"$polyapi_hide$,omitempty"`
}

// HTTPResponse is the response of http request
type HTTPResponse struct {
	StatusCode int
	Status     string
	Body       json.RawMessage
	Header     Header
}

// DoRequestAPI is the custom api for access apis from polyapi
func (c *PolyClient) DoRequestAPI(apiPath string, method string, header Header, body interface{}) (*HTTPResponse, error) {
	bodyBytes, err := c.GenHeaderSignature(header, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPRequest(c.remoteURL+apiPath, method, header, bodyBytes)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	r := &HTTPResponse{
		Body:       respBody,
		Header:     resp.Header,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	return r, nil
}

// CustomBody is custom body
type CustomBody map[string]interface{}

// Add insert a new field to custom body
func (b CustomBody) Add(name string, data interface{}) bool {
	return b.add(name, data, false)
}

// Set insert a new field to custom body forced
func (b CustomBody) Set(name string, data interface{}) bool {
	return b.add(name, data, true)
}

// // SetSignature set bodySignature for custome body
// func (b CustomBody) SetSignature(c *PolyClient) bool {
// 	return b.add(BodySinature, c.bodySign.genBodySignature(), true)
// }

func (b CustomBody) add(name string, data interface{}, force bool) bool {
	if _, ok := b[name]; !ok && !force {
		return false
	}
	b[name] = data
	return true
}

// NewCustomBody generate a custom body with signature
func (c *PolyClient) NewCustomBody() CustomBody {
	return CustomBody{
		//BodySinature: c.GenBodySignature(),
	}
}

// MakeBodyBase create a BodyBase with signature
func (c *PolyClient) MakeBodyBase() BodyBase {
	return BodyBase{
		//Signature: c.GenBodySignature(),
	}
}

// // GenBodySignature generate body signature
// func (c *PolyClient) GenBodySignature() json.RawMessage {
// 	return json.RawMessage(c.bodySign.genBodySignature())
// }

// HTTPRequest do a custom http request
func (c *PolyClient) HTTPRequest(reqURL, method string, header Header, data []byte) (*http.Response, error) {
	if err := validateHTTPMethod(method); err != nil {
		return nil, err
	}

	uri, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if method == MethodGet { // data to query if method is 'GET'
		uri.RawQuery = signature.ToQuery(data)
	} else {
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, uri.String(), reader)
	if err != nil {
		return nil, err
	}
	req.Header = header

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// GenHeaderSignature create header.Signature form body and return body bytes
func (c *PolyClient) GenHeaderSignature(header Header, body interface{}) ([]byte, error) {
	var b []byte
	var err error
	switch d := body.(type) {
	case []byte:
		b = d
	default:
		if b, err = json.Marshal(body); err != nil {
			return nil, err
		}
	}

	signInfo := polysign.PolySignatureInfo{
		AccessKeyID: c.accessKeyID,
		SignMethod:  polysign.XHeaderPolySignMethodVal,
		SignVersion: polysign.XHeaderPolySignVersionVal,
		Timestamp:   timestamp(),
		Body:        b,
	}

	signature, err := c.sign.Signature(&signInfo)
	if err != nil {
		return nil, err
	}
	header.Set(polysign.XHeaderPolySignVersion, signInfo.SignVersion)
	header.Set(polysign.XHeaderPolySignMethod, signInfo.SignMethod)
	header.Set(polysign.XHeaderPolySignKeyID, signInfo.AccessKeyID)
	header.Set(polysign.XHeaderPolySignTimestamp, signInfo.Timestamp)

	header.Set(polysign.XHeaderPolySignSignature, signature)

	return b, nil
}

func validateHTTPMethod(method string) error {
	switch method {
	case MethodGet, MethodHead, MethodPost,
		MethodPut, MethodPatch, MethodDelete,
		MethodConnect, MethodOptions, MethodTrace:
		return nil
	}
	return fmt.Errorf("unsupport http method '%s'", method)
}
