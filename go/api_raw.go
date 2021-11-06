package polysdk

import (
	"net/http"
)

// RequestRawAPI request raw api from polyapi
// apiPath is the full namespace of raw api. eg: /system/raw/sample_raw_api
func (c *PolyClient) RequestRawAPI(apiPath string, method string, header http.Header, data []byte) error {
	return nil
}

// RawAPIDoc request raw api from polyapi
// apiPath is the full namespace of raw api. eg: /system/raw/sample_raw_api
func (c *PolyClient) RawAPIDoc(apiPath string, method string, header http.Header, data []byte) error {
	return nil
}
