package gormutil

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConnectWithDSN(t *testing.T) {
	viper.SetEnvPrefix("DATABASE")
	_ = viper.BindEnv("DSN")
	_ = viper.Unmarshal(&C)

	db, err := ConnectWithDSN(C.DSN, C)
	assert.NoError(t, err)
	assert.IsType(t, &gorm.DB{}, db)
}
