package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"os"
	"polysdk/consts"
	"polysdk/internal/config"
	"polysdk/internal/crypto/aesx"
	"strings"
)

const hashMethod = "sha256"

var stringEncoder = base64.RawURLEncoding

// Signer is the interface for Signature
type Signer interface {
	Signature(data interface{}) (string, error)
}

// NewSigner create a hmac hash generator.
// NOTE: It accept crypted secret key only.
func NewSigner(cryptedSecretKey string) (Signer, error) {
	sk, err := config.Decrypt(cryptedSecretKey)
	if err != nil {
		return nil, err
	}
	return newSigner(sk), nil
}

// NewSignerFromFile create a hmac hash generator from config file.
func NewSignerFromFile(cfgFilePath string) (Signer, *config.PolyConfig, error) {
	if cfgFilePath == "" { // try to use env var
		cfgFilePath = os.Getenv(consts.EnvPolyConfigPath)
	}

	cfg, err := config.LoadFromFile(cfgFilePath, false)
	if err != nil {
		return nil, nil, err
	}

	keys, err := cfg.GetCryptoKeys()
	if err != nil {
		return nil, nil, err
	}
	sk, err := aesx.DecodeString(cfg.Key.SecretKey, keys...)
	if err != nil {
		return nil, nil, err
	}

	return newSigner(sk), cfg.HideSecret(), nil
}

func newSigner(secretKey string) *hmacHash {
	return newHmac([]byte(secretKey), hashMethod)
}

// NewHmac create a hmac hash
func newHmac(secretKey []byte, hashMethod string) *hmacHash {
	var h hash.Hash
	switch strings.ToLower(hashMethod) {
	case "sha256":
		h = hmac.New(sha256.New, secretKey)
	default:
		panic("unsupport hash method:" + hashMethod)
	}
	return &hmacHash{h: h}
}

// the hmac hash
type hmacHash struct {
	h hash.Hash
}

// Signature generate signature for data
func (h *hmacHash) Signature(data interface{}) (string, error) {
	query := ToQuery(data)
	// println(query)
	// println(_toRawQuery(data))

	signatureBytes, err := h.hash([]byte(query))
	if err != nil {
		return "", err
	}

	signature := stringEncoder.EncodeToString(signatureBytes)
	return signature, nil
}

// generate hmac hash
func (h *hmacHash) hash(entity []byte) ([]byte, error) {
	h.h.Reset() // reset the hmac state
	return sha(entity, h.h)
}

func sha(entity []byte, hash hash.Hash) ([]byte, error) {
	_, err := hash.Write(entity)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
