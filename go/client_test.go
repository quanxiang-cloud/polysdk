package polysdk_test

import (
	"polysdk"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := polysdk.NewPolyClient("")
	if err != nil {
		panic(err)
	}
	h := polysdk.Header{}
	h.Set("Content-Type", "application/json")
	body := map[string]interface{}{
		"time_stamp": polysdk.Timestamp(""),
		"zone":       "pek3d",
		"_signature": c.GenBodySignature(),
		"active":     -1,
	}
	polysdk.PrettyShow(body)

	//uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM3"
	uri := "/api/v1/polyapi/namespace/tree/system/app/jhdsk"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}
	polysdk.PrettyShow(r)

}
