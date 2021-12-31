package test_tool

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geeksmy/go-libs/jwt"
)

func JwtSetup(t *testing.T, testJwtSecret string) {
	err := jwt.SetupWithConf(jwt.Conf{Secret: testJwtSecret})
	assert.NoError(t, err)
}
