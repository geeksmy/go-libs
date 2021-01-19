package jwt

import (
	"crypto/rsa"
	"io/ioutil"
	"os"

	jwtlib "github.com/dgrijalva/jwt-go"
)

// CommonTokenConfig is common token-related configuration settings.
// The setting are used by TokenProvider and TokenValidator.
type CommonTokenConfig struct {
	TokenName   string `json:"token_name,omitempty" xml:"token_name" yaml:"token_name"`
	TokenSecret string `json:"token_secret,omitempty" xml:"token_secret" yaml:"token_secret"`
	TokenIssuer string `json:"token_issuer,omitempty" xml:"token_issuer" yaml:"token_issuer"`
	TokenOrigin string `json:"token_origin,omitempty" xml:"token_origin" yaml:"token_issuer"`
	// The expiration time of a token in seconds
	TokenLifetime      int    `json:"token_lifetime,omitempty" xml:"token_lifetime" yaml:"token_lifetime"`
	TokenSigningMethod string `json:"token_signing_method,omitempty" xml:"token_signing_method" yaml:"token_signing_method"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ParseRSAPrivateKeyFromPEM(key string) (*rsa.PrivateKey, error) {
	// 检查 key 是不是文件
	// 如果是则从文件获取 key
	if fileExists(key) {
		content, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		key = string(content)
	}

	return jwtlib.ParseRSAPrivateKeyFromPEM([]byte(key))
}

func ParseRSAPublicKeyFromPEM(key string) (*rsa.PublicKey, error) {
	// 检查 key 是不是文件
	// 如果是则从文件获取 key
	if fileExists(key) {
		content, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		key = string(content)
	}

	return jwtlib.ParseRSAPublicKeyFromPEM([]byte(key))
}
