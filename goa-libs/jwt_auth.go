package goalibs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/geeksmy/go-libs/jwt"
	"github.com/geeksmy/go-libs/util"
	"goa.design/goa/v3/security"
)

type CtxKey int

const (
	CurrentUserIDKey CtxKey = 999
	JwtClaimsKey     CtxKey = 1000
)

var (
	ErrorUnauthorized = errors.New("请登录后再试")
)

type JwtAuth struct {
}

// JWT 认证
func (j *JwtAuth) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	// parse && verify JWT token,
	validator := jwt.NewValidator()

	// 1. parse JWT token,
	userClaims, err := validator.Verify(token, scheme)
	if err != nil {
		return ctx, err
	}

	// 2. validate provided "scopes" claim
	if err := validateScopes(scheme.RequiredScopes, userClaims.Scopes); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, CurrentUserIDKey, userClaims.ID)
	ctx = context.WithValue(ctx, JwtClaimsKey, userClaims)

	return ctx, nil
}

// 获取当前登录用户ID
func (j *JwtAuth) GetCurrentUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(CurrentUserIDKey).(string)
	if !ok {
		return "", ErrorUnauthorized
	}
	return userID, nil
}

// 范围验证，支持使用通配符匹配
// 比如 user:*, *:read, *, *:delete
func validateScopes(expected, actual []string) error {
	var missing []string
	for _, r := range expected {
		found := false
		for _, s := range actual {
			// if s == r {
			if util.Glob(s, r) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, r)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing scopes: %s", strings.Join(missing, ", "))
}
