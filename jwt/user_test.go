package jwt

import (
	"reflect"
	"testing"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// testPriKey2 is the same as "test-priv-2.key"
var userTestRSPriKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgEMFBKcGW7iRRlJdIuF0/5YmB3ACsCd6hWCFk4FGAj7G+sd4m9GG
U/9ae9x00yvkY2Pit03B5kxHQfVAqKG6PnTzRg5cbwjPjnhFiPeLfGWMKIIEkhTa
cuIu8Tr+hmMchxCUYl9twakFl3bOVsHqmMcByJ44FII66Kl4z6k4ERKZAgMBAAEC
gYAfGugi4SeWzQ43UfTLcTLirDnNeeHqIMpglv50BFssacug4tBm+ZJotMVB95K/
D1w10tbCpxjNFFF/k4fwr/EmeuAK3aQgmsbxAgtH6hyKtYp6yrK7jabkXXJLFTaC
8aWgq7RRCazDxlJlOtn50vMUH1LHf1Z0YUC76OyzsiKC9QJBAINN8Nl11M4/3s1n
x4H0sMiyyW8DhqMrpla0IgAwuWRHmWZ1VuiWUXmv/oW+YLoFxDofukhLFT2NblFr
h5d4kW8CQQCCqnoG2Wd0fRFk1kHcGEZzJB0D1PKepOHe//ca4uNPupo45qOXaMCU
7vj7+JkZo/pEgjXaG1G00saF5KTMJgh3AkA+F82eCKrqHiou2LTwL9aqEmJPrUsu
PqYaunSZwnDpizJv0W2X7/33ndKvTKhRUAjLs9VT+q3AvfE9b6xfZRThAkBVifKe
fz45xRJY9+ZfhkjAYbjY5FP8RSZUjS6gHD4A2MDTVTFtEjdYiGTY1vKrFWzl4nQM
l2vSu1UZHAhCWPebAkAT9KpSzWqcLt7GFOHjoVpHIeuyCCkWJwS9JeP6J/QbaJq/
SMNiwTaDC1kT8uCWqTgd5u5AKOV+oyzwmj0nJu8n
-----END RSA PRIVATE KEY-----`

func TestUserHSClaims(t *testing.T) {
	claims := &UserClaims{}
	claims.ExpiresAt = time.Now().Add(time.Duration(900) * time.Second).Unix()
	claims.Name = "Greenberg, Paul"
	claims.Email = "greenpau@outlook.com"
	claims.Origin = "localhost"
	claims.Subject = "greenpau@outlook.com"
	claims.Roles = append(claims.Roles, "anonymous")
	secret := "75f03764-147c-4d87-b2f0-4fda89e331c8"
	token, err := claims.GetToken(SigningMethodHS512, []byte(secret))
	if err != nil {
		t.Fatalf("Failed to get JWT token for %v: %s", claims, err)
	}
	t.Logf("Token: %s", token)
}

func TestAppMetadataAuthorizationRoles(t *testing.T) {
	secret := "75f03764147c4d87b2f04fda89e331c808ab50a932914e758ae17c7847ef27fa"
	encodedToken := "eyJhbGciOiJIUzUxMiJ9.eyJyb2xlcyI6WyJhZG1pbiIsImVkaXRvciIsImd1ZXN0Il19" +
		".Ynd1E4BASvR5E5BfGKx-8-pAuwEuaMNT2UzR9wO8taYiWEadFDAR38PevAeZCSOMY6QfF5DXvQWhotin6X7aQw"
	expectedRoles := []string{"admin", "editor", "guest"}

	t.Logf("token Secret: %s", secret)
	t.Logf("encoded Token: %s", encodedToken)

	token, err := jwtlib.Parse(encodedToken, func(token *jwtlib.Token) (interface{}, error) {
		if _, validMethod := token.Method.(*jwtlib.SigningMethodHMAC); !validMethod {
			return nil, ErrUnexpectedSigningMethod.WithArgs(token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed parsing the encoded token: %s", err)
	}

	t.Logf("token: %v", token)

	claimMap := token.Claims.(jwtlib.MapClaims)
	claims, err := NewUserClaimsFromMap(claimMap)
	if err != nil {
		t.Fatalf("failed parsing claims for token: %s", err)
	}

	t.Logf("claims: %v", claims)

	if len(claims.Roles) == 0 {
		t.Fatalf("no roles found, expecting %s", expectedRoles)
	}

	if len(claims.Roles) != len(expectedRoles) {
		t.Fatalf("role count mismatch: %d (token) vs %d (expected)", len(claims.Roles), len(expectedRoles))
	}

	if !reflect.DeepEqual(claims.Roles, expectedRoles) {
		t.Fatalf("role mismatch: %s (token) vs %s (expected)", claims.Roles, expectedRoles)
	}

	t.Logf("token roles: %s", claims.Roles)
}

func TestAnonymousGuestRoles(t *testing.T) {
	secret := "75f03764147c4d87b2f04fda89e331c808ab50a932914e758ae17c7847ef27fa"
	encodedToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
		"eyJleHAiOjI1NDI3MTkzOTgsInN1YiI6ImdyZWVucGF1QG91dGxvb2suY29tIiwibmFtZSI6IkdyZW" +
		"VuYmVyZywgUGF1bCIsImVtYWlsIjoiZ3JlZW5wYXVAb3V0bG9vay5jb20iLCJvcmlnaW4iOiJsb2Nh" +
		"bGhvc3QifQ." +
		"INRBEsx7b4sewCmNCQxRSN3Hk_sT5BMbjlq_iPdbvkYiWnORS93xYSAei78GWEMDepc6ECTSGhqVL-sDFCbPoA"
	expectedRoles := []string{"anonymous", "guest"}

	t.Logf("token Secret: %s", secret)
	t.Logf("encoded Token: %s", encodedToken)

	token, err := jwtlib.Parse(encodedToken, func(token *jwtlib.Token) (interface{}, error) {
		if _, validMethod := token.Method.(*jwtlib.SigningMethodHMAC); !validMethod {
			return nil, ErrUnexpectedSigningMethod.WithArgs(token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed parsing the encoded token: %s", err)
	}

	t.Logf("token: %v", token)

	claimMap := token.Claims.(jwtlib.MapClaims)
	claims, err := NewUserClaimsFromMap(claimMap)
	if err != nil {
		t.Fatalf("failed parsing claims for token: %s", err)
	}

	t.Logf("claims: %v", claims)

	if len(claims.Roles) == 0 {
		t.Fatalf("no roles found, expecting %s", expectedRoles)
	}

	if len(claims.Roles) != len(expectedRoles) {
		t.Fatalf("role count mismatch: %d (token) vs %d (expected)", len(claims.Roles), len(expectedRoles))
	}

	if !reflect.DeepEqual(claims.Roles, expectedRoles) {
		t.Fatalf("role mismatch: %s (token) vs %s (expected)", claims.Roles, expectedRoles)
	}

	t.Logf("token roles: %s", claims.Roles)
}

func TestUserRSClaims(t *testing.T) {
	claims := &UserClaims{}
	claims.ExpiresAt = time.Now().Add(time.Duration(900) * time.Second).Unix()
	claims.Name = "Jones, Nika"
	claims.Email = "njones@outlook.example.com"
	claims.Origin = "localhost"
	claims.Subject = "njones@outlook.example.com"
	claims.Roles = append(claims.Roles, "anonymous")
	claims.MetaData = map[string]string{"gender": "male"}

	priKey, err := jwtlib.ParseRSAPrivateKeyFromPEM([]byte(userTestRSPriKey))
	if err != nil {
		t.Fatal(err)
	}
	token, err := claims.GetToken(SigningMethodRS512, priKey)
	if err != nil {
		t.Fatalf("Failed to get JWT token for %v: %s", claims, err)
	}
	t.Logf("Token: %s", token)

	validator, _ := NewValidatorImplWithConf(Conf{
		RSAPrivateKey: userTestRSPriKey,
	})

	newClaims, err := validator.Verify(token, nil)
	assert.NoError(t, err)
	assert.Equal(t, newClaims.MetaData["gender"], claims.MetaData["gender"])
}
