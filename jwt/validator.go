package jwt

import (
	"crypto/rsa"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
)

// Validator Errors
const (
	ErrNoBackends       strError = "no token backends available"
	ErrExpiredToken     strError = "expired token"
	ErrNoAccessList     strError = "user role is valid, but denied by default deny on empty access list"
	ErrAccessNotAllowed strError = "user role is valid, but not allowed by access list"
	ErrNoParsedClaims   strError = "failed to extract claims"
	ErrNoTokenFound     strError = "no token found"

	ErrInvalidParsedClaims strError = "failed to extract claims: %s"
	ErrInvalidSecret       strError = "secret key backend error: %s"
	ErrInvalid             strError = "%v"
)

// TokenValidator validates tokens in http requests.
type TokenValidator struct {
	CommonTokenConfig
	Cache         *TokenCache
	AccessList    []*AccessListEntry
	TokenBackends []TokenBackend

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewTokenValidator returns an instance of TokenValidator
func NewTokenValidator() *TokenValidator {
	v := &TokenValidator{}

	v.Cache = NewTokenCache()
	v.TokenLifetime = 900
	return v
}

// ConfigureTokenBackends configures available TokenBackend.
func (v *TokenValidator) ConfigureTokenBackends() error {
	v.TokenBackends = []TokenBackend{}
	if v.PrivateKey != nil || v.PublicKey != nil {
		backend := NewRSAKeyTokenBackend(v.PrivateKey, v.PublicKey)
		v.TokenBackends = append(v.TokenBackends, backend)
	}
	if v.TokenSecret != "" {
		backend, err := NewSecretKeyTokenBackend(v.TokenSecret)
		if err != nil {
			return ErrInvalidSecret.WithArgs(err)
		}
		v.TokenBackends = append(v.TokenBackends, backend)
	}
	if len(v.TokenBackends) == 0 {
		return ErrNoBackends
	}
	return nil
}

// ValidateToken parses a token and returns claims, if valid.
func (v *TokenValidator) ValidateToken(token string) (*UserClaims, bool, error) {
	valid := false
	// First, check cached entries
	claims := v.Cache.Get(token)
	if claims != nil {
		if claims.ExpiresAt < time.Now().Unix() {
			_ = v.Cache.Delete(token)
			return nil, false, ErrExpiredToken
		}
		valid = true
	}

	parseErrors := []error{}
	// If not valid, parse claims from a string.
	if !valid {
		for _, backend := range v.TokenBackends {
			mapClaims := jwtlib.MapClaims{}

			token, err := jwtlib.ParseWithClaims(token, &mapClaims, backend.ProvideKey)
			if err != nil {
				parseErrors = append(parseErrors, err)
				continue
			}
			if !token.Valid {
				continue
			}
			if mapClaims == nil {
				parseErrors = append(parseErrors, ErrInvalid.WithArgs("claims is nil"))
				continue
			}

			claims, err = NewUserClaimsFromMap(mapClaims)
			if err != nil {
				parseErrors = append(parseErrors, ErrInvalidParsedClaims.WithArgs(err))
				continue
			}

			valid = true
			break
		}
	}

	if valid {
		if len(v.AccessList) == 0 {
			return nil, false, ErrNoAccessList
		}
		aclAllowed := false
		for _, entry := range v.AccessList {
			claimAllowed, abortProcessing := entry.IsClaimAllowed(claims)
			if abortProcessing {
				aclAllowed = claimAllowed
				break
			}

			if claimAllowed {
				aclAllowed = true
			}
		}
		if !aclAllowed {
			return nil, false, ErrAccessNotAllowed
		}
	}

	if !valid {
		// return nil, false, ErrInvalid.WithArgs(errorMessages)
		return nil, false, parseErrors[0]
	}

	return claims, true, nil
}
