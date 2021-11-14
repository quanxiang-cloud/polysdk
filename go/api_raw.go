package polysdk

import (
	"polysdk/internal/apipath"
)

// doc type
const (
	DocRaw        = "raw"
	DocSwag       = "swag"
	DocCurl       = "curl"
	DocJavascript = "javascript"
	DocPython     = "python"
)

// RawAPIRequest request raw api from polyapi
// apiPath is the full namespace of raw api. eg: /system/raw/sample_raw_api
func (c *PolyClient) RawAPIRequest(fullNamespace string, method string, header Header, data interface{}) (*HTTPResponse, error) {
	uri := apipath.Join(apipath.APIRawRequest, fullNamespace)
	return c.DoRequestAPI(uri, method, header, data)
}

type apiDocReq struct {
	BodyBase
	DocType    string `json:"docType"`
	TitleFirst bool   `json:"titleFirst"`
}

// RawAPIDoc request raw api from polyapi
// fullNamespace is the full namespace of raw api. eg: /system/raw/sample_raw_api
func (c *PolyClient) RawAPIDoc(fullNamespace string, docType string, titleFirst bool) (*HTTPResponse, error) {
	d := apiDocReq{
		BodyBase: BodyBase{
			Signature: c.bodySign.genBodySignature(),
		},
		DocType:    docType,
		TitleFirst: titleFirst,
	}
	header := Header{
		HeaderContentType: []string{MIMEJSON},
	}
	uri := apipath.Join(apipath.APIRawDoc, fullNamespace)
	return c.DoRequestAPI(uri, MethodPost, header, d)
}
