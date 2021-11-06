package polysdk

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	x, _ := newSignatureTemplate("$access_key_id$")
	fmt.Println(x.timestampIndex, x.signatureTemplate)
	fmt.Println(string(x.genBodySignature()))
}
