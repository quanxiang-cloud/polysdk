package polysdk

import (
	"polysdk/internal/apipath"
)

// PolyAPIRequest request poly api from polyapi
// fullNamespace is the full namespace of poly api. eg: /system/poly/sample_poly_api
func (c *PolyClient) PolyAPIRequest(fullNamespace string, method string, header Header, data interface{}) (*HttpResponse, error) {
	uri := apipath.Join(apipath.APIPolyRequest, fullNamespace)
	return c.DoRequestAPI(uri, method, header, data)
}

// PolyAPIDoc request poly api from polyapi
// fullNamespace is the full namespace of poly api. eg: /system/poly/sample_poly_api
func (c *PolyClient) PolyAPIDoc(fullNamespace string, docType string, titleFirst bool) (*HttpResponse, error) {
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
	uri := apipath.Join(apipath.APIPolyDoc, fullNamespace)
	return c.DoRequestAPI(uri, MethodPost, header, d)
}
