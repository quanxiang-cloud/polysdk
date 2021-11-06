package polysdk

import (
	"polysdk/internal/signature"
)

// NewPolyClient create a ploy client from config file
func NewPolyClient(configPath string) (*PolyClient, error) {
	sign, cfg, err := signature.NewSignerFromFile(configPath)
	if err != nil {
		return nil, err
	}
	bodySign, err := newSignatureTemplate(cfg.Key.AccessKeyID)
	if err != nil {
		return nil, err
	}
	return &PolyClient{
		remoteURL: cfg.RemoteURL,
		sign:      sign,
		bodySign:  bodySign,
	}, nil
}

// PolyClient is a client for polyapi
type PolyClient struct {
	remoteURL string
	sign      signature.Signer
	bodySign  *signatureTemplate
}
