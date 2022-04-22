package polysdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"polysdk/internal/hash"
	"polysdk/internal/polysign"
	"polysdk/internal/signature"
)

// header define
const (
	HeaderContentType = "Content-Type"
	HeaderXRequestID  = "X-Request_id"
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
	// NOTE: path and other hide parameter
	PolyHide map[string]interface{} `json:"$polyapi_hide$,omitempty"`

	// NOTE: none-object customer body root
	CustomerBody interface{} `json:"$body$,omitempty"`
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

	var respBody []byte
	if resp.StatusCode == http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		respBody = b
	} else {
		respBody = []byte(fmt.Sprintf(`%q`, resp.Status))
	}

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

func (b CustomBody) add(name string, data interface{}, force bool) bool {
	if _, ok := b[name]; ok && !force {
		return false
	}
	b[name] = data
	return true
}

// NewCustomBody generate a custom body with signature
func (c *PolyClient) NewCustomBody() CustomBody {
	return CustomBody{}
}

// MakeBodyBase create a BodyBase with signature
func (c *PolyClient) MakeBodyBase() BodyBase {
	return BodyBase{}
}

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
		uri.RawQuery, err = signature.ToQuery(data)
		if err != nil {
			return nil, err
		}
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
	case json.RawMessage:
		b = []byte(d)
	case []byte:
		b = d
	default:
		if b, err = json.Marshal(body); err != nil {
			return nil, err
		}
	}

	header.Set(HeaderXRequestID, hash.ShortID(0))

	var signInfo CustomBody
	if err := json.Unmarshal(b, &signInfo); err != nil {
		return nil, err
	}

	var (
		signVals = [][2]string{
			{polysign.XHeaderPolySignKeyID, c.accessKeyID},
			{polysign.XHeaderPolySignMethod, polysign.XHeaderPolySignMethodVal},
			{polysign.XHeaderPolySignVersion, polysign.XHeaderPolySignVersionVal},
			{polysign.XHeaderPolySignTimestamp, c.polyTimestamp()},
		}
		opSignVals = func(fn func(string, string) error) error {
			for _, v := range signVals {
				if err := fn(v[0], v[1]); err != nil {
					return err
				}
			}
			return nil
		}
		opAddBodyHeader = func(name, val string) error {
			if name != polysign.XBodyPolySignSignature {
				header.Set(name, val)
			}
			if !signInfo.Add(name, val) {
				err = fmt.Errorf("duplicate field %s in body", name)
				return err
			}
			return nil
		}
		opDelBody = func(name, val string) error {
			delete(signInfo, name)
			return nil
		}
	)

	if err := opSignVals(opAddBodyHeader); err != nil {
		return nil, err
	}

	signature, err := c.sign.Signature(signInfo)
	if err != nil {
		return nil, err
	}
	if err := opAddBodyHeader(polysign.XBodyPolySignSignature, signature); err != nil {
		return nil, err
	}

	if err := opSignVals(opDelBody); err != nil {
		return nil, err
	}

	return json.Marshal(signInfo)
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
