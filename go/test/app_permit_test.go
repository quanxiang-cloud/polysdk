package polysdk_test

import (
	"fmt"
	"io"
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

func _TestHttpRequest(t *testing.T) {
	body := map[string]interface{}{
		"appID":  "jhdsk",
		"userID": "893ca81d-f571-4a6f-8088-673e8775ff64",
	}
	//没有该分组访问权限 893ca81d-f571-4a6f-8088-673e8775ff64 jhdsk false false

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, polysdk.MIMEJSON)

	//uri := "http://app-center/api/v1/app-center/checkIsAdmin"
	uri := "http://app-center/api/v1/app-center/checkAppAccess"
	bodyBytes, err := c.GenHeaderSignature(h, body)
	if err != nil {
		panic(err)
	}

	resp, err := c.HTTPRequest(uri, polysdk.MethodPost, h, bodyBytes)
	if err != nil {
		panic(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	println(uri, resp.Status)
	fmt.Println("reqBody", string(bodyBytes))
	fmt.Println("respHeader", resp.Header)
	fmt.Println("respBody", string(respBody))
}

func _TestHttpRequest2(t *testing.T) {
	body := map[string]interface{}{
		"appID":    "jhdsk",
		"userID":   "893ca81d-f571-4a6f-8088-673e8775ff64",
		"is_super": false,
	}

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, polysdk.MIMEJSON)

	uri := "http://app-center/api/v1/app-center/checkIsAdmin"
	bodyBytes, err := c.GenHeaderSignature(h, body)
	if err != nil {
		panic(err)
	}

	resp, err := c.HTTPRequest(uri, polysdk.MethodPost, h, bodyBytes)
	if err != nil {
		panic(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	println(uri, resp.Status)
	fmt.Println("reqBody", string(bodyBytes))
	fmt.Println("respHeader", resp.Header)
	fmt.Println("respBody", string(respBody))
}

func _TestHttpRequest3(t *testing.T) {
	body := map[string]interface{}{
		"appID":    "jhdsk",
		"userID":   "893ca81d-f571-4a6f-8088-673e8775ff64",
		"is_super": false,
	}

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, polysdk.MIMEJSON)

	uri := "http://polyapi_inner:9090/api/v1/polyapi/inner/regSwaggerAlone/system/app/jhdsk"
	bodyBytes, err := c.GenHeaderSignature(h, body)
	if err != nil {
		panic(err)
	}

	resp, err := c.HTTPRequest(uri, polysdk.MethodPost, h, bodyBytes)
	if err != nil {
		panic(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	println(uri, resp.Status)
	fmt.Println("reqBody", string(bodyBytes))
	fmt.Println("respHeader", resp.Header)
	fmt.Println("respBody", string(respBody))
}

func _TestClient2(t *testing.T) {
	body := map[string]interface{}{
		"time_stamp": polysdk.Timestamp(""),
		"zone":       "pek3d",
		"_signature": c.GenBodySignature(),
		"active":     -1,
		"name":       "test3",
		"appID":      "cfb879cf-157b-4d5d-8200-a4ae350209ff",
		"pathType":   "raw.3party",
	}
	polysdk.PrettyShow(body)

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	//uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM3"
	//uri := "/api/v1/polyapi/namespace/tree/system/app/jhdsk"
	uri := "/api/v1/polyapi/namespace/appPath"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}

func TestClient3(t *testing.T) {
	body := map[string]interface{}{
		"time_stamp": polysdk.Timestamp(""),
		"zone":       "pek3d",
		"_signature": c.GenBodySignature(),
		"active":     -1,
		"name":       "test3",
		"appID":      "cfb879cf-157b-4d5d-8200-a4ae350209ff",
		"pathType":   "raw.3party",
		"key":        "3c84556041d79f354dc20085332d28ce",
		"city":       "110101",
	}
	polysdk.PrettyShow(body)

	h := polysdk.Header{}
	h.Set(polysdk.HeaderContentType, "application/json")

	//uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM3"
	//uri := "/api/v1/polyapi/namespace/tree/system/app/jhdsk"
	//uri := "/api/v1/polyapi/poly/request/system/app/bwp2w/poly/test01/ceshi1"
	//uri := "/api/v1/polyapi/poly/request/system/app/fqjqv/poly/demo/test"
	//uri := "/api/v1/polyapi/poly/request/system/app/rtftw/poly/baron01/linxiao001/kxd01"
	//uri := "/api/v1/polyapi/poly/query/system/app/fqjqv/poly/demo/test"
	uri := "/api/v1/polyapi/poly/list/system/app/28hr9/poly/p"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)
}
