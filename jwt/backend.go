package jwt

import (
	"crypto/rsa"

	jwtlib "github.com/dgrijalva/jwt-go"
)

// Backend Errors
const (
	ErrInvalidSecretLength strError = "secrets less than 16 characters in length are not allowed"
	ErrUnexpectedKID       strError = "the kid specified in the header was not found"
	ErrNoRSAKeyFound       strError = "no RSA key found"

	ErrUnexpectedSigningMethod strError = "signing method mismatch: %v (expected) vs. %v (received)"
)

// TokenBackend is the interface to provide key material.
type TokenBackend interface {
	ProvideKey(token *jwtlib.Token) (interface{}, error)
}

// SecretKeyTokenBackend hold symentric keys from HS family.
type SecretKeyTokenBackend struct {
	secret []byte
}

// NewSecretKeyTokenBackend returns SecretKeyTokenBackend instance.
func NewSecretKeyTokenBackend(s string) (*SecretKeyTokenBackend, error) {
	if len(s) < 16 {
		return nil, ErrInvalidSecretLength
	}
	b := &SecretKeyTokenBackend{
		secret: []byte(s),
	}
	return b, nil
}

// ProvideKey provides key material from SecretKeyTokenBackend.
func (b *SecretKeyTokenBackend) ProvideKey(token *jwtlib.Token) (interface{}, error) {
	if _, validMethod := token.Method.(*jwtlib.SigningMethodHMAC); !validMethod {
		return nil, ErrUnexpectedSigningMethod.WithArgs("HS", token.Header["alg"])
	}
	return b.secret, nil
}

// RSAKeyTokenBackend hold asymentric keys from RS family.
type RSAKeyTokenBackend struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAKeyTokenBackend returns RSKeyTokenBackend instance.
func NewRSAKeyTokenBackend(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSAKeyTokenBackend {
	b := &RSAKeyTokenBackend{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	return b
}

// ProvideKey provides key material from RSKeyTokenBackend.
func (b *RSAKeyTokenBackend) ProvideKey(token *jwtlib.Token) (interface{}, error) {
	if _, validMethod := token.Method.(*jwtlib.SigningMethodRSA); !validMethod {
		return nil, ErrUnexpectedSigningMethod.WithArgs("RS", token.Header["alg"])
	}

	if b.privateKey != nil {
		return b.privateKey.Public(), nil
	} else if b.publicKey != nil {
		return b.publicKey, nil
	}

	return nil, ErrNoRSAKeyFound
}
