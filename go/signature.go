package polysdk

import (
	"encoding/json"
	"fmt"
	"polysdk/internal/hash"
	"polysdk/internal/polysign"
	"sync/atomic"
	"time"
)

// Now return time that adjust with server.
func (c *PolyClient) Now() time.Time {
	t := time.Now()
	if adj := c.getTimeAdjust(); adj != 0 {
		t = t.Add(adj)
	}
	return t
}

// SetTimeAdjustMS change local time adjust value to coordinate with server.
// NOTE: It never change local system clock.
func (c *PolyClient) SetTimeAdjustMS(adjust int64) {
	timeAdjust := time.Duration(adjust) * time.Millisecond
	c.setTimeAdjust(timeAdjust)
}

func (c *PolyClient) getTimeAdjust() time.Duration {
	return time.Duration(atomic.LoadInt64(&c.timeAdjust))
}

func (c *PolyClient) setTimeAdjust(adjust time.Duration) {
	fmt.Printf("SetTimeAdjustMS %s => %s\n", c.getTimeAdjust(), adjust)

	timeAdjust := int64(adjust)
	atomic.StoreInt64(&c.timeAdjust, timeAdjust)
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

func systemNow() time.Time {
	return time.Now()
}

// SyncServerClock auto adjust time with server
func (c *PolyClient) SyncServerClock() error {
	uri := "/api/v1/gate/ping"
	body := polysign.PingReq{
		Rand:          hash.ShortID(0),
		PingTimestamp: systemNow().Format(polysign.PingTimestampFmt),
	}
	h := Header{}
	h.Set(HeaderContentType, "application/json")
	r, err := c.DoRequestAPI(uri, MethodPost, h, body)
	if err != nil {
		return err
	}

	var respBody struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		Data polysign.PingResp `json:"data"`
	}
	var d *polysign.PingResp = &respBody.Data
	if err := json.Unmarshal(r.Body, &respBody); err != nil {
		return err
	}
	ping, err := time.Parse(polysign.PingTimestampFmt, d.PingTimestamp)
	if err != nil {
		return err
	}
	pong, err := time.Parse(polysign.PingTimestampFmt, d.PongTimestamp)
	if err != nil {
		return err
	}

	now := systemNow()
	netDelay := now.Sub(ping) / 10
	adjust := pong.Add(netDelay).Sub(now)
	c.setTimeAdjust(adjust)

	// if true {
	// 	println("adjust", adjust.String())
	// 	println("netDelay", netDelay.String())
	// 	println("ping", d.PingTimestamp, ping.String())
	// 	println("pong", d.PongTimestamp, pong.String())
	// 	println("now", now.String())
	// }

	return nil
}
