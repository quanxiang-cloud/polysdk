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
		"active":     -1,
		"name":       "test3",
	}
	polysdk.PrettyShow(body)

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	//uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM3"
	//uri := "/api/v1/polyapi/namespace/tree/system/app/jhdsk"
	uri := "/api/v1/polyapi/namespace/create/system/app/swhnm/poly"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func TestRawRequest(t *testing.T) {
	body := polysdk.CustomBody{
		//"time_stamp": polysdk.Timestamp(""),
		//"zone":       "pek3d",
		polysdk.XPolyBodyHideArgs: map[string]interface{}{
			"app": "appX",
		},
		polysdk.XPolyCustomerBodyRoot: "foo",
	}

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	//uri := "/system/app/jhdsk/customer/ns2/viewVM3"
	//uri := "/system/app/2zfvz/customer/three_group/timeTicker"
	//uri := "/system/app/gskgx/raw/customer/a/bb.r"
	uri := "/system/app/gskgx/raw/inner/form/form/2q2bh_get.r"
	r, err := c.RawAPIRequest(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func _TestRawDoc(t *testing.T) {
	apiPath := "/system/app/gskgx/raw/customer/a/bb.r"
	r, err := c.RawAPIDoc(apiPath, polysdk.DocJavascript, true)
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
