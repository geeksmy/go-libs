package jwt

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	goa "goa.design/goa/v3/pkg"
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
	rsaKeyPair2 = [2]string{
		`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCsFMXBSc0NuvD6X9uF9QoHgqpPUsoQGqjb/EVevT+PfQaJJpiG
pPmjZYeouu0keBVBryLaPsHMLnqPkmo8O9okrm6nKd3j8dh1D1J/gGn3uXbFIVL8
xBou4d4IjrQoUQy1AkDSZqp5D1SNH5tBptQHykfyZWjW1nz5tWVGNTrv9wIDAQAB
AoGBAJpCAMh1nQDK7Qy083gRIo16/6seHw7ypx6U9aao5Zp+TGaUF7TTnQlxkXub
HcxMbVm1YvvbuCSOIcYkTWFzkedJnrEaYnKaKLNY8dmYDPE39omRVf1fKunAM0H7
7wKjjex91CEj2kY8mRkXGuWzyr94bDlfvgEo31FBYzmHBJPhAkEA10FsjiUcp95F
jLSdhDKt3FnWIa7wQzVsyoxe9wI3n1kr9rBJAj7zsZDdc3FJIBSMxaQL6LgMBMXx
a8dQ0B2NBwJBAMynQhb+YHmYTyFDtdFcnSde36bm3kuv4CqMd2bcQOwf3y+L2jzI
PWpaEMz/d+gh2tOBk7hwdDJS0C1TbqjMuZECQEFYTBMOsxdGw4hGYDb8h4kIAJgz
Gh7c/gyy9jU4CBiook7+DvvOjn4OAxwDfPZvJpjtBux7yrI8QOC+Hgs/nUsCQQC6
dREz3eOMJYbC6ev3qgfM3RWr/OA/2SfF3GDVKItGFuvDaAoYIuBBy3lPLNrUdjPn
TEGIY3yadPUSttc1mADhAkAzDcIG7Dql4tFVVHhLDCUAJM/Zq4jeuOFS3ItFEyvs
jfXML6f/ekIaH+5AusfJdHFUuijE/OwblsptVMrqvAp5
-----END RSA PRIVATE KEY-----`,
		`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsFMXBSc0NuvD6X9uF9QoHgqpP
UsoQGqjb/EVevT+PfQaJJpiGpPmjZYeouu0keBVBryLaPsHMLnqPkmo8O9okrm6n
Kd3j8dh1D1J/gGn3uXbFIVL8xBou4d4IjrQoUQy1AkDSZqp5D1SNH5tBptQHykfy
ZWjW1nz5tWVGNTrv9wIDAQAB
-----END PUBLIC KEY-----`,
	}
	rsaKeyPair3 = [2]string{
		"tests/jwt.key",
		"tests/jwt.key.pub",
	}
)

func TestSignerImpl_Sign(t *testing.T) {
	_, err := NewSignerImplWithConf(Conf{})
	assert.Equal(t, err, ErrNoRequiredSecret)

	cases := []struct {
		caseName string
		conf     Conf
		method   SigningMethod
		claim    jwt.Claims
		err      error
	}{
		{
			"测试不支持的 Claims",
			Conf{RSAPrivateKey: rsaKeyPair1[0], RSAPublicKey: rsaKeyPair1[1]},
			SigningMethodRS512,
			jwt.MapClaims{},
			ErrNotSupportedClaims,
		},

		{
			"测试使用 RSA 签名",
			Conf{RSAPrivateKey: rsaKeyPair1[0], RSAPublicKey: rsaKeyPair1[1]},
			SigningMethodRS512,
			UserClaims{ID: "123456"},
			nil,
		},

		{
			"测试使用 HMAC 签名",
			Conf{Secret: "ulGoc6DKy4Ur3i+xAuOGQSS4Q3AJcCuEBzTgRkm3WSM"},
			SigningMethodHS512,
			UserClaims{ID: "123456"},
			nil,
		},

		{
			"测试使用 RSA 签名(通过文件)",
			Conf{RSAPrivateKey: rsaKeyPair3[0], RSAPublicKey: rsaKeyPair3[1]},
			SigningMethodRS256,
			UserClaims{ID: "123456"},
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("case: %s", tc.caseName)

		signer, err := NewSignerImplWithConf(tc.conf)
		assert.NoError(t, err)

		token, err := signer.Sign(tc.claim)
		if tc.err == nil {
			assert.NoError(t, err)

			validator, err1 := NewValidatorImplWithConf(tc.conf)
			assert.NoError(t, err1)

			_, err = validator.Verify(token, nil)
			assert.NoError(t, err)
		} else {
			assert.Equal(t, tc.err, err)
		}
	}
}

func signRSATokenStr(priKey string, claims jwt.Claims) (string, error) {
	privateKey, err := ParseRSAPrivateKeyFromPEM(priKey)
	if err != nil {
		return "", err
	}

	if claims == nil {
		claims = jwt.MapClaims{}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(privateKey)
}

func TestJwtValidatorImpl_Verify_RSA(t *testing.T) {
	type testCase struct {
		keyPair [2]string
		err     error
		claim   UserClaims
		token   string
	}
	cases := []testCase{
		// test pass
		{rsaKeyPair1, nil, UserClaims{ID: "123456"}, "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIxMjM0NTYifQ.O69ZHHKofbvtG0yv5In14rhsH" +
			"_2ptRGN2hxB--aLv6P1UYC26zX7v9_-lUR4vTY7Wuwutv_6JjIMWTp4DuyKp3Xlx73tyoL-9B8iEbMMd0zsfPxdBBYw8gFTUp" +
			"-zHOMFKNU8F6r4ug1pgnIvVP7uK9ApTK1XABiA2LetF8MJ9aQ"},
		// test rsa key with file
		{
			rsaKeyPair3,
			goa.TemporaryError(
				"unauthorized",
				translateJwtValidationError(jwt.NewValidationError("err for UT", jwt.ValidationErrorSignatureInvalid)),
			),
			UserClaims{ID: "123456"},
			"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIxMjM0NTYifQ.O69ZHHKofbvtG0yv5In14rhsH" +
				"_2ptRGN2hxB--aLv6P1UYC26zX7v9_-lUR4vTY7Wuwutv_6JjIMWTp4DuyKp3Xlx73tyoL-9B8iEbMMd0zsfPxdBBYw8gFTUp" +
				"-zHOMFKNU8F6r4ug1pgnIvVP7uK9ApTK1XABiA2LetF8MJ9aQ"},
		// signature error
		{
			[2]string{rsaKeyPair2[0], rsaKeyPair2[1]},
			// goa.TemporaryError("unauthorized", jwtValidationErrDisplay[jwt.ValidationErrorSignatureInvalid]),
			goa.TemporaryError(
				"unauthorized",
				translateJwtValidationError(jwt.NewValidationError("err for UT", jwt.ValidationErrorSignatureInvalid)),
			),
			UserClaims{},
			"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.e30.EmEyxWgqqs7p_JKZnzK6kFWE0b2A6z_ORNmfWjHt1cN9gdmH8PxE6kUgdiNFbswwxtMLxnHzRyYtT3N" +
				"LFubeafvAlJ3TNNgWUnaVfcfyAHCrXz74ZTYhbI_artbl4s0AMR4o-q92jCIGlaSOO-U_U1tiZ4UYeL1Ua-mbs1kMS0k",
		},
		// expired
		{
			rsaKeyPair1,
			goa.TemporaryError(
				"unauthorized",
				translateJwtValidationError(jwt.NewValidationError("err for UT", jwt.ValidationErrorExpired)),
			),
			UserClaims{ExpiresAt: time.Now().AddDate(0, 0, -1).Unix()},
			"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTg5NTgwMjh9.UYyqlv04JcVZ_YK1Tbyamw_8D3Udbdn1OtgWuBqPvLbdCX__fgilllzO3Wv" +
				"tO7r396F0WWZJHiPtf6neynRqSw4BZ16XiOkoTCjRSmH7KFlZ8vQ8lEKU4_ohR2KW3fUEHa9WTf-B_5x2JFLKKJSj-RyivTBwunAbhMNoWoAJyJI",
		},
	}

	for idx, tc := range cases {
		conf := Conf{
			RSAPrivateKey: tc.keyPair[0],
			RSAPublicKey:  tc.keyPair[1],
		}

		var err error
		sharedOption, err = newJwtOption(conf)
		assert.NoError(t, err)
		if err != nil {
			t.Logf("case %d failed", idx)
			continue
		}

		_, err1 := signRSATokenStr(tc.keyPair[0], tc.claim)
		assert.NoError(t, err1)
		if err1 != nil {
			t.Logf("case %d failed", idx)
			continue
		}

		validator := NewValidatorImpl()
		claims, err := validator.Verify(tc.token, nil)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
			continue
		}

		assert.Equal(t, claims.ID, tc.claim.ID)
	}
}
