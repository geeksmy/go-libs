package jwt

import (
	"errors"
	"fmt"

	jwtlib "github.com/dgrijalva/jwt-go"
)

type strError string

func (e strError) Error() string { return string(e) }

// F captures the values for an error string formatting.
func (e strError) WithArgs(v ...interface{}) error {
	var hasErr, hasNil bool
	for _, vv := range v {
		switch err := vv.(type) {
		case error:
			if err == nil {
				return nil // we pass nill errors along
			}
			hasErr = true
		case nil:
			hasNil = true
		}
	}

	if hasNil && !hasErr {
		return nil
	}

	return fmtErr{err: fmt.Errorf("%w", e), v: v}
}

// fmtErr is for errors that will be formatted. It holds
// formatting values in a slice so they can be added when the
// error is stringfied. Otherwise the underlining error without
// formatting can be matched.
type fmtErr struct {
	err error
	v   []interface{}
}

func (e fmtErr) Error() string { return fmt.Sprintf(e.err.Error(), e.v...) }

// Unwrap is a method to help unwrap errors on the base error for go1.13+
func (e fmtErr) Unwrap() error { return errors.Unwrap(e.err) }

func translateJwtValidationError(err *jwtlib.ValidationError) string {
	var jwtValidationErrDisplay = map[uint32]string{
		// Signature validation failed
		jwtlib.ValidationErrorMalformed: "令牌格式错误",
		// Token could not be verified because of signing problems
		jwtlib.ValidationErrorUnverifiable: "令牌签名错误无法验证",
		// Signature validation failed
		jwtlib.ValidationErrorSignatureInvalid: "签名验证失败",

		// Standard Claim validation errors
		// AUD validation failed
		jwtlib.ValidationErrorAudience: "受众验证失败",
		// EXP validation failed
		jwtlib.ValidationErrorExpired: "令牌已过期",
		// IAT validation failed
		jwtlib.ValidationErrorIssuedAt: "发布时间验证失败",
		// ISS validation failed
		jwtlib.ValidationErrorIssuer: "申请人验证失败",
		// NBF validation failed
		jwtlib.ValidationErrorNotValidYet: "令牌尚未生效",
		// JTI validation failed
		jwtlib.ValidationErrorId: "JTI 验证失败",
		// Generic claims validation error
		jwtlib.ValidationErrorClaimsInvalid: "通用声明验证错误",
	}

	s, ok := jwtValidationErrDisplay[err.Errors]
	if ok {
		return s
	}

	return "令牌错误"
}
