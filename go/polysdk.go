package polysdk

import (
	"net"
	"net/http"
	"polysdk/internal/signature"
	"time"
)

const (
	httpTimeout      = 10
	httpMaxIdleConns = 3
)

// NewPolyClient create a ploy client from config file
func NewPolyClient(configPath string) (*PolyClient, error) {
	sign, cfg, err := signature.NewSignerFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return &PolyClient{
		remoteURL:   cfg.RemoteURL,
		accessKeyID: cfg.Key.AccessKeyID,
		sign:        sign,
		httpClient: http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(time.Second * httpTimeout)
					c, err := net.DialTimeout(netw, addr, time.Second*httpTimeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
				MaxIdleConns:      httpMaxIdleConns,
				DisableKeepAlives: false,
			},
		},
	}, nil
}

// PolyClient is a client for polyapi
type PolyClient struct {
	remoteURL   string
	accessKeyID string
	sign        signature.Signer
	httpClient  http.Client
}
