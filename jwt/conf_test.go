package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfErr_WithMessage(t *testing.T) {
	err := ConfErr{
		code:    ConfErrCodeRSAPrivateKey,
		message: "original message",
	}

	e := err.WithMessage("cover by new message")
	assert.Equal(t, e.message, "original message: \"cover by new message\"")
	assert.Equal(t, err.message, "original message")
}

func TestConf_Validate(t *testing.T) {
	type testCase struct {
		err  error
		conf Conf
	}

	cases := []testCase{
		{err: ErrNoRequiredSecret, conf: Conf{}},
		{err: ErrInvaliRSAPrivateKey, conf: Conf{
			RSAPrivateKey: "invalid rsa private key",
		}},
		{err: ErrInvalidRSAPublicKey, conf: Conf{
			RSAPublicKey: "invalid rsa public key",
		}},
	}

	for _, tc := range cases {
		_, err := newJwtOption(tc.conf)
		if tc.err == nil {
			continue
		}

		// nolint(errorlint): fixme
		expect := tc.err.(ConfErr)
		// nolint(errorlint): fixme
		e := err.(ConfErr)
		assert.Equal(t, expect.code, e.Code())
	}
}
