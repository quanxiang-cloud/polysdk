package polysdk

import (
	"polysdk/internal/polysign"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	println(time.Now().Format(polysign.PingTimestampFmt))
}
