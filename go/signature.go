package polysdk

import (
	"polysdk/internal/polysign"
	"time"
)

// Timestamp get current timestamp of UTC
func Timestamp(f string) string {
	if f == "" {
		f = "2006-01-02T15:04:05Z"
	}
	return time.Now().UTC().Format(f)
}

func timestamp() string {
	return time.Now().Format(polysign.XHeaderPolySignTimestampFmt)
}
