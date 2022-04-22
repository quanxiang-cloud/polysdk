package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"polysdk/consts"
	"polysdk/internal/crypto/aesx"
	"polysdk/internal/crypto/deviceid"
	"polysdk/internal/hash"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// NewInitConfig create a default config object
func NewInitConfig(desc string) *PolyConfig {
	return &PolyConfig{
		RemoteURL:   consts.DefaultRemoteURL,
		CreateAt:    time.Now().Format("2006-01-02T15:04:05MST"),
		Description: desc,
	}
}

// LoadFromFile create a config object from file
func LoadFromFile(filePath string, validate bool) (*PolyConfig, error) {
	p := &PolyConfig{}
	if err := p.LoadFile(filePath); err != nil {
		return nil, err
	}
	if validate {
		if err := p.Validate(); err != nil {
			return nil, fmt.Errorf("%s %s", filePath, err.Error())
		}
	}
	return p, nil

}

// PolyKeyConfig is the config of poly apikey
type PolyKeyConfig struct {
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
}

// Empty check if key is empty
func (k PolyKeyConfig) Empty() bool {
	return k.AccessKeyID == "" && k.SecretKey == ""
}

// PolyConfig is the config of
type PolyConfig struct {
	RemoteURL   string        `json:"remoteUrl"`
	Key         PolyKeyConfig `json:"key"`
	CreateAt    string        `json:"createAt"`
	Description string        `json:"description"`

	enc encoding
}

// LoadFile load poly config from file
func (c *PolyConfig) LoadFile(filePath string) error {
	if err := c.setEncoding(filePath); err != nil {
		return err
	}

	b, err := LoadFile(filePath)
	if err != nil {
		return err
	}

	return c.unmarshal(b)
}

// StoreFile store poly config into file
func (c *PolyConfig) StoreFile(filePath string, allowOverwrite bool) error {
	if !c.Key.Empty() { // validate non empty key
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if err := c.setEncoding(filePath); err != nil {
		return err
	}
	b, err := c.marshal()
	if err != nil {
		return err
	}

	return StoreFile(b, filePath, allowOverwrite)
}

// GetCryptoKeys get crypt keys for aesx
func (c *PolyConfig) GetCryptoKeys() ([]string, error) {
	if c.RemoteURL == "" || c.Key.AccessKeyID == "" {
		return nil, fmt.Errorf("missing $.remoteUrl or $.key.accessKeyId")
	}

	keys := []string{
		c.RemoteURL,
		c.Key.AccessKeyID,
		c.CreateAt,
		deviceid.MyHexDeviceID(),
	}
	return []string{
		hash.HexString(hash.Sha256Hash(nil, 97, 65535, keys...)),
	}, nil
}

// Validate verify if this config is ready for signature
func (c *PolyConfig) Validate() error {
	if c.Key.SecretKey == "" {
		return fmt.Errorf("missing $.key.secretKey")
	}

	keys, err := c.GetCryptoKeys()
	if err != nil {
		return err
	}
	if _, err := aesx.DecodeString(c.Key.SecretKey, keys...); err != nil {
		return fmt.Errorf("poly config validate fail")
	}
	return nil
}

// HideSecret clean secret key
func (c *PolyConfig) HideSecret() *PolyConfig {
	c.Key.SecretKey = ""
	return c
}

//------------------------------------------------------------------------------

// LoadFile open and read all bytes from file
func LoadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	return b, err
}

// StoreFile store bytes to overwrite file
func StoreFile(data []byte, filePath string, allowOverwrite bool) error {
	flag := os.O_CREATE | os.O_WRONLY
	if allowOverwrite {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}

	f, err := os.OpenFile(filePath, flag, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

// Encrypt encryt plant text using strong aes method
func Encrypt(plantTxt string) string {
	key, err := aesx.EncodeString(plantTxt, deviceid.MyHexDeviceID())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return key
}

// Decrypt decrypt plant text by strong aes method
func Decrypt(cryptTxt string) (string, error) {
	key, err := aesx.DecodeString(cryptTxt, deviceid.MyHexDeviceID())
	if err != nil {
		return "", err
	}
	return key, nil
}

//------------------------------------------------------------------------------

func (c *PolyConfig) marshal() ([]byte, error) {
	switch c.enc {
	case jsonBinding:
		return json.MarshalIndent(c, "", "  ")
	case yamlBinding:
		return yaml.Marshal(c)
	default:
		return nil, fmt.Errorf("unsupport encoding %v", c.enc)
	}
}

func (c *PolyConfig) unmarshal(data []byte) error {
	switch c.enc {
	case jsonBinding:
		return json.Unmarshal(data, c)
	case yamlBinding:
		return yaml.Unmarshal(data, c)
	default:
		return fmt.Errorf("unsupport encoding %v", c.enc)
	}
}

func (c *PolyConfig) setEncoding(filePath string) error {
	enc, err := c.parseEncoding(filePath)
	if err != nil {
		return err
	}
	c.enc = enc
	return nil
}

// parse encoding format from file extension
func (c *PolyConfig) parseEncoding(filePath string) (encoding, error) {
	switch ext := strings.ToLower(filepath.Ext(filePath)); ext {
	case ".json", ".jsn":
		return jsonBinding, nil
	case ".yaml", ".yml":
		return yamlBinding, nil
	default:
		return 0, fmt.Errorf("unsupport file extension [%s]", ext)
	}
}

type encoding uint8

// encoding format
const (
	jsonBinding encoding = iota + 1
	yamlBinding
)
