package test_tool

import (
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	libsRedis "github.com/geeksmy/go-libs/redis"
)

func RedisCnnForTest(t *testing.T, envURIKey string) *redis.Client {
	viper.AutomaticEnv()
	uri := viper.GetString(envURIKey)

	c := libsRedis.Conf{
		URI: uri,
	}

	connector, err := libsRedis.NewConnector(c)
	assert.NoError(t, err)
	cli, err := connector.Connect()
	assert.NoError(t, err)

	libsRedis.Client = cli
	return cli

}
