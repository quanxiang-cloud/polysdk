package polysdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"polysdk/consts"
	"polysdk/internal/signature"
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

// BodyBase is base struct of body
type BodyBase struct {
	// {"version":1,"method":"HmacSHA256","access_key_id":"$access_key_id$","timestamp":"2006-01-02T15:04:05-0700"}
	Signature json.RawMessage `json:"_signature"`
	// path and other parameter
	Hide interface{} `json:"_hide,omitempty"`
}

// HttpResponse is the response of http request
type HttpResponse struct {
	StatusCode int
	Status     string
	Body       json.RawMessage
	Header     Header
}

func (c *PolyClient) DoRequestAPI(apiPath string, method string, header Header, body interface{}) (*HttpResponse, error) {
	bodyBytes, err := c.genHeaderSignature(header, body)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest(c.remoteURL+apiPath, method, header, bodyBytes)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	r := &HttpResponse{
		Body:       respBody,
		Header:     resp.Header,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	return r, nil
}

func HttpRequest(reqURL, method string, header Header, data []byte) (*http.Response, error) {
	if err := validateHttpMethod(method); err != nil {
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

	resp, err := http.DefaultClient.Do(req) // http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//------------------------------------------------------------------------------

func (c *PolyClient) GenBodySignature() json.RawMessage {
	return json.RawMessage(c.bodySign.genBodySignature())
}

func (c *PolyClient) genHeaderSignature(header Header, body interface{}) ([]byte, error) {
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
	signature, err := c.sign.Signature(b)
	if err != nil {
		return nil, err
	}
	header.Set(consts.HeaderSignature, signature)
	return b, nil
}

func validateHttpMethod(method string) error {
	switch method {
	case MethodGet, MethodHead, MethodPost,
		MethodPut, MethodPatch, MethodDelete,
		MethodConnect, MethodOptions, MethodTrace:
		return nil
	}
	return fmt.Errorf("unsupport http method '%s'", method)
}
