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

func TestHttpRequest(t *testing.T) {
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

func TestHttpRequest2(t *testing.T) {
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

func TestHttpRequest3(t *testing.T) {
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
