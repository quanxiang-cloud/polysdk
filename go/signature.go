package polysdk

import (
	"polysdk/internal/polysign"
	"time"
)

// Now return time that adjust with server.
func (c *PolyClient) Now() time.Time {
	t := time.Now()
	if c.timeAdjust != 0 {
		t = t.Add(c.timeAdjust)
	}
	return t
}

// SetTimeAdjustMS change local time adjust value to coordinate with server.
// NOTE: It never change local system clock.
func (c *PolyClient) SetTimeAdjustMS(adjust int) {
	c.timeAdjust = time.Duration(adjust) * time.Millisecond
}

// Timestamp get current timestamp of UTC
func (c *PolyClient) Timestamp(f string) string {
	if f == "" {
		f = "2006-01-02T15:04:05Z"
	}
	return c.Now().UTC().Format(f)
}

func (c *PolyClient) polyTimestamp() string {
	return c.Now().Format(polysign.XHeaderPolySignTimestampFmt)
}

// SyncServerClock auto adjust time with server
func (c *PolyClient) SyncServerClock() error {
	// TODO
	return nil
}
