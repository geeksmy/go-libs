package jwt

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
)

// User Errors
const (
	ErrInvalidClaimExpiresAt strError = "invalid exp type"
	ErrInvalidClaimIssuedAt  strError = "invalid iat type"
	ErrInvalidClaimNotBefore strError = "invalid nbf type"
	ErrInvalidSigningMethod  strError = "unsupported signing method"
	ErrUnsupportedSecret     strError = "empty secrets are not supported"

	ErrInvalidRole                strError = "invalid role type %T in roles"
	ErrInvalidRoleType            strError = "invalid roles type %T"
	ErrInvalidScopesType          strError = "invalid scopes type %T"
	ErrInvalidOrg                 strError = "invalid org type %T in orgs"
	ErrInvalidOrgType             strError = "invalid orgs type %T"
	ErrInvalidAppMetadataRoleType strError = "invalid roles type %T in app_metadata-authorization"

	ErrInvalidConfiguration        strError = "%s: default access list configuration error: %s"
	ErrInvalidBackendConfiguration strError = "%s: token validator configuration error: %s"
)

// UserClaims represents custom and standard JWT claims.
// https://tools.ietf.org/html/rfc7519#section-4.1
type UserClaims struct {
	Audience      string            `json:"aud,omitempty" xml:"aud" yaml:"aud,omitempty"`
	ExpiresAt     int64             `json:"exp,omitempty" xml:"exp" yaml:"exp,omitempty"`
	ID            string            `json:"jti,omitempty" xml:"jti" yaml:"jti,omitempty"`
	IssuedAt      int64             `json:"iat,omitempty" xml:"iat" yaml:"iat,omitempty"`
	Issuer        string            `json:"iss,omitempty" xml:"iss" yaml:"iss,omitempty"`
	NotBefore     int64             `json:"nbf,omitempty" xml:"nbf" yaml:"nbf,omitempty"`
	Subject       string            `json:"sub,omitempty" xml:"sub" yaml:"sub,omitempty"`
	Name          string            `json:"name,omitempty" xml:"name" yaml:"name,omitempty"`
	Email         string            `json:"email,omitempty" xml:"email" yaml:"email,omitempty"`
	Roles         []string          `json:"roles,omitempty" xml:"roles" yaml:"roles,omitempty"`
	Origin        string            `json:"origin,omitempty" xml:"origin" yaml:"origin,omitempty"`
	Scopes        []string          `json:"scopes,omitempty" xml:"scopes" yaml:"scopes,omitempty"`
	Organizations []string          `json:"org,omitempty" xml:"org" yaml:"org,omitempty"`
	MetaData      map[string]string `json:"meta,omitempty" xml:"meta" yaml:"meta,omitempty"`
}

// Valid validates user claims.
func (u UserClaims) Valid() error {
	if u.ExpiresAt < time.Now().Unix() {
		return errors.New("token expired")
	}
	return nil
}

// NewUserClaimsFromMap returns UserClaims.
func NewUserClaimsFromMap(m map[string]interface{}) (*UserClaims, error) {
	u := &UserClaims{}

	if _, exists := m["aud"]; exists {
		u.Audience = m["aud"].(string)
	}
	if _, exists := m["exp"]; exists {
		switch exp := m["exp"].(type) {
		case float64:
			u.ExpiresAt = int64(exp)
		case json.Number:
			v, _ := exp.Int64()
			u.ExpiresAt = v
		default:
			return nil, ErrInvalidClaimExpiresAt
		}
	}

	if _, exists := m["jti"]; exists {
		u.ID = m["jti"].(string)
	}

	if _, exists := m["iat"]; exists {
		switch exp := m["iat"].(type) {
		case float64:
			u.IssuedAt = int64(exp)
		case json.Number:
			v, _ := exp.Int64()
			u.IssuedAt = v
		default:
			return nil, ErrInvalidClaimIssuedAt
		}
	}

	if _, exists := m["iss"]; exists {
		u.Issuer = m["iss"].(string)
	}

	if _, exists := m["nbf"]; exists {
		switch exp := m["nbf"].(type) {
		case float64:
			u.NotBefore = int64(exp)
		case json.Number:
			v, _ := exp.Int64()
			u.NotBefore = v
		default:
			return nil, ErrInvalidClaimNotBefore
		}
	}

	if _, exists := m["sub"]; exists {
		u.Subject = m["sub"].(string)
	}

	if _, exists := m["name"]; exists {
		u.Name = m["name"].(string)
	}

	if _, exists := m["mail"]; exists {
		u.Email = m["mail"].(string)
	}

	if _, exists := m["email"]; exists {
		u.Email = m["email"].(string)
	}

	if _, exists := m["roles"]; exists {
		switch m["roles"].(type) {
		case []interface{}:
			roles := m["roles"].([]interface{})
			for _, role := range roles {
				switch role := role.(type) {
				case string:
					u.Roles = append(u.Roles, role)
				default:
					return nil, ErrInvalidRole.WithArgs(role)
				}
			}
		case string:
			roles := m["roles"].(string)
			u.Roles = append(u.Roles, strings.Split(roles, " ")...)
		default:
			return nil, ErrInvalidRoleType.WithArgs(m["roles"])
		}
	}

	if _, exists := m["origin"]; exists {
		u.Origin = m["origin"].(string)
	}

	if _, exists := m["scopes"]; exists {
		switch m["scopes"].(type) {
		case []interface{}:
			scopes := m["scopes"].([]interface{})
			for _, scope := range scopes {
				switch scope := scope.(type) {
				case string:
					u.Scopes = append(u.Scopes, scope)
				default:
					return nil, ErrInvalidScopesType.WithArgs(scope)
				}
			}
		case string:
			scopes := m["scopes"].(string)
			u.Scopes = append(u.Scopes, strings.Split(scopes, " ")...)
		default:
			return nil, ErrInvalidScopesType.WithArgs(m["scopes"])
		}
	}

	if _, exists := m["org"]; exists {
		switch m["org"].(type) {
		case []interface{}:
			orgs := m["org"].([]interface{})
			for _, org := range orgs {
				switch org := org.(type) {
				case string:
					u.Organizations = append(u.Organizations, org)
				default:
					return nil, ErrInvalidOrg.WithArgs(org)
				}
			}
		case string:
			orgs := m["org"].(string)
			u.Organizations = append(u.Organizations, strings.Split(orgs, " ")...)
		default:
			return nil, ErrInvalidOrgType.WithArgs(m["org"])
		}
	}

	u.MetaData = make(map[string]string)
	if _, exists := m["meta"]; exists {
		switch m["meta"].(type) {
		case map[string]interface{}:
			meta := m["meta"].(map[string]interface{})
			for k, v := range meta {
				switch v := v.(type) {
				case string:
					u.MetaData[k] = v
				default:
					return nil, ErrInvalidOrg.WithArgs(v)
				}
			}
		default:
			return nil, ErrInvalidOrgType.WithArgs(m["meta"])
		}
	}

	if len(u.Roles) == 0 {
		u.Roles = append(u.Roles, "anonymous")
		u.Roles = append(u.Roles, "guest")
	}

	return u, nil
}

func (u *UserClaims) FromJSON(data []byte) error {
	return json.Unmarshal(data, u)
}

// GetToken returns a signed JWT token
func (u *UserClaims) GetToken(method SigningMethod, secret interface{}) (string, error) {
	if secret == nil {
		return "", ErrUnsupportedSecret
	}

	sm := method.getSigningMethod()

	token := jwtlib.NewWithClaims(sm, u)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
