package goalibs

import (
	"context"
	"testing"
	"time"

	"git.chinaopen.ai/yottacloud/go-libs/jwt"

	"github.com/stretchr/testify/assert"
	"goa.design/goa/v3/security"
)

var (
	rsaKeyPair1 = [2]string{
		`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDStzmdY4WVtawqmh8iMfYO+ZG0MdJ4hG2bX78YdRYQuQeD8Qt5
gvz32XVIBK6oWll4jQ+kWcRPePLCqYVp0D0fa/DyrxUDZMYrVPgAfK4hLLuqNUSj
J9XMTqQf3wouYSETbop53wzIkzNt9v0UxzjGsKd2eXrrCu/XL4hIDD2TewIDAQAB
AoGAJKJaT/S3itm0/wsgko9hGiVloZBv6SgM7lBtJtFkbq/ckKHdvth5JpYV/9lg
jEB5Aa50o7w/lxmOCy3x1f2wQm9C1n6p0/1omYT1fE+lsW7hBqCsmgLdJPQEtKU/
hknSS2ka+2AZMNhqPKDIPbdO8dx3h6HZ19imZZjOV9qurCECQQD6ZLs5x36RefBS
Sxn1YeHVbhOhc/xpP5M4oKWLAo/i3SpU+BSZGKWn6ha9UJtYqGppKxvRSMnH8qsk
sLxxFNGJAkEA128PBu7S/HIitcR2nNJyYBmnQV/7R+4TrFT6CYWswKrgQtBw+i3Q
HGwNgelC70pecwrX/0sfAPj8yaQrMgrP4wJAecGZuRLk3OCql5mlWcjrA9mu+gdo
adVcRWD4nov0Y6CgAQlRGtD9XJVb0Vrs70OV2uX3dOEsHWpHBmRafbJd4QJAZstY
ST+Vrt0rEL5Ato14A8DZ3MjWlrjT7EvTzs1Y48nlHhxlfAp2Hcm1HELYLGO4832+
GNOYZJdpzzElZ0hfbwJAP3SJZXHG7m5ffHF3QBPcYBblZtn2Qekc0hNPUtR8wei9
HXbImuSMKfUvEjpUDnm57UV3hgWST26hm+tT56cBeA==
-----END RSA PRIVATE KEY-----

`,
		`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDStzmdY4WVtawqmh8iMfYO+ZG0
MdJ4hG2bX78YdRYQuQeD8Qt5gvz32XVIBK6oWll4jQ+kWcRPePLCqYVp0D0fa/Dy
rxUDZMYrVPgAfK4hLLuqNUSjJ9XMTqQf3wouYSETbop53wzIkzNt9v0UxzjGsKd2
eXrrCu/XL4hIDD2TewIDAQAB
-----END PUBLIC KEY-----
`,
	}

	userID = "09ce10ee-fe90-4ce5-967c-e46cc209fb58"
	scopes = []string{"api:write", "api:read"}
)

func TestSignerImpl_Valid(t *testing.T) {
	setupJwt()

	claims := newValidClaims()
	signer := jwt.NewSigner()

	token, err1 := signer.Sign(claims)
	assert.NoError(t, err1)

	auther := JwtAuth{}

	_, err2 := auther.JWTAuth(context.Background(), token, newScheme())
	assert.NoError(t, err2)
}

func TestSignerImpl_No_Scopes(t *testing.T) {
	setupJwt()

	claims := newNoScopesClaims()
	signer := jwt.NewSigner()

	token, err1 := signer.Sign(claims)
	assert.NoError(t, err1)

	auther := JwtAuth{}

	_, err2 := auther.JWTAuth(context.Background(), token, newScheme())
	assert.Error(t, err2)
}

func TestSignerImpl_Expired(t *testing.T) {
	setupJwt()

	claims := newExpiredClaims()
	signer := jwt.NewSigner()

	token, err1 := signer.Sign(claims)
	assert.NoError(t, err1)

	auther := JwtAuth{}

	_, err2 := auther.JWTAuth(context.Background(), token, newScheme())
	assert.Error(t, err2)
}

func setupJwt() {
	jwt.C.RSAPrivateKey = rsaKeyPair1[0]
	jwt.C.RSAPublicKey = rsaKeyPair1[1]
	_ = jwt.Init()
}

func newValidClaims() jwt.UserClaims {
	return jwt.UserClaims{
		ID:        userID,
		Scopes:    scopes,
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
	}
}

func newNoScopesClaims() jwt.UserClaims {
	return jwt.UserClaims{
		ID: userID,
		// Scopes:     scopes,
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
	}
}

func newExpiredClaims() jwt.UserClaims {
	return jwt.UserClaims{
		ID: userID,
		// Scope:    scopes,
		ExpiresAt: time.Now().AddDate(0, 0, -1).Unix(),
	}
}

func newScheme() *security.JWTScheme {
	return &security.JWTScheme{
		Name:           "jwt",
		Scopes:         []string{"api:read", "api:write"},
		RequiredScopes: []string{"api:read"},
	}
}

func Test_ValidateScopes(t *testing.T) {
	err1 := validateScopes([]string{"user:read"}, []string{"user:*"})
	assert.NoError(t, err1)

	err2 := validateScopes([]string{"user:read", "user:write"}, []string{"user:read"})
	assert.EqualError(t, err2, "missing scopes: user:write")

	err3 := validateScopes([]string{"user:read", "user:write"}, []string{"*"})
	assert.NoError(t, err3)

	err4 := validateScopes([]string{"user:read"}, []string{"*:read"})
	assert.NoError(t, err4)

	err5 := validateScopes([]string{"user:read"}, []string{"user:read"})
	assert.NoError(t, err5)

	err6 := validateScopes([]string{"user:read"}, []string{"group:read", "orgin:read"})
	assert.EqualError(t, err6, "missing scopes: user:read")

	err7 := validateScopes([]string{"user:read"}, []string{"group:read", "user:read"})
	assert.NoError(t, err7)
}
