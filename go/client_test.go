package polysdk_test

import (
	"polysdk"
	"testing"
)

var c *polysdk.PolyClient

func init() {
	_c, err := polysdk.NewPolyClient("")
	if err != nil {
		panic(err)
	}
	c = _c
}

func _TestClient(t *testing.T) {
	body := map[string]interface{}{
		"time_stamp": polysdk.Timestamp(""),
		"zone":       "pek3d",
		"_signature": c.GenBodySignature(),
		"active":     -1,
	}
	polysdk.PrettyShow(body)

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	//uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM3"
	uri := "/api/v1/polyapi/namespace/tree/system/app/jhdsk"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func _TestRawRequest(t *testing.T) {
	body := polysdk.CustomBody{
		"time_stamp":         polysdk.Timestamp(""),
		"zone":               "pek3d",
		polysdk.BodySinature: c.GenBodySignature(),
	}

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	uri := "/system/app/jhdsk/customer/ns2/viewVM3"
	r, err := c.RawAPIRequest(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func _TestRawDoc(t *testing.T) {
	apiPath := "/system/form/base_pergroup_create"
	r, err := c.RawAPIDoc(apiPath, polysdk.DocSwag, false)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func _TestPolyRequest(t *testing.T) {
	body := polysdk.CustomBody{
		"appID":       "app1",
		"name":        "app1Name",
		"description": "description",
		"scopes": []polysdk.CustomBody{
			polysdk.CustomBody{
				"type": 1,
				"id":   "someid1",
				"name": "somename1",
			},
			polysdk.CustomBody{
				"type": 2,
				"id":   "someid2",
				"name": "somename2",
			},
		},
		polysdk.BodySinature: c.GenBodySignature(),
	}
	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")
	uri := "/system/poly/permissionInit"
	r, err := c.PolyAPIRequest(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func _TestPolyDoc(t *testing.T) {
	apiPath := "/system/poly/permissionInit"
	r, err := c.PolyAPIDoc(apiPath, polysdk.DocRaw, false)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}
