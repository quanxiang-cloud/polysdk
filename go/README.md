# polysdk-go

polysdk-go is a client for polyapi.

## usage

1. set environment value of ENV_POLY_CONFIG_PATH
It will share config file from tool and client.
```
 set ENV_POLY_CONFIG_PATH=c:\data\poly_cfg.json
```
2. initial an empty config file
```
polykit i
```
3. edit the config file with AccessKeyID and Secret
```
{
  "remoteUrl": "http://polyapi.qxp.com",
  "key": {
    "accessKeyId": "<key-id>",
    "secretKey": "<secret-key>"
  },
  "createAt": "2021-11-12T06:35:03CST",
  "description": ""
}
```
4. encrypt the config file
```
polykit c
```
5. verify the encrpted config file
```
polykit v
```
6. create a poly client from the config file
```Go
	c, err := polysdk.NewPolyClient("")
	if err != nil {
		panic(err)
	}
	c.SyncServerClock() // adjust local clock with server
	
	h := polysdk.Header{}
	h.Set("Content-Type", "application/json")
	body := map[string]interface{}{
		"zone":       "pek3d",
	}
	polysdk.PrettyShow(body)

	uri := "/api/v1/polyapi/raw/request/system/app/jhdsk/customer/ns2/viewVM.r"
	r, err := c.DoRequestAPI(uri, polysdk.MethodPost, h, body)
	if err != nil {
		panic(err)
	}

	polysdk.PrettyShow(r)
```

~~ enjoy this sdk!