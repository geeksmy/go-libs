package gormutil

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
# default yaml file
Database:
    DSN: "postgres://postgres:postgres@localhost/test"
    LogMode: false
    MaxIdleConns: 0
    MaxOpenConns: 0
    ConnMaxLifetime: 0
*/
// Conf 定义 db 配置的 struct
type Conf struct {
	DSN             string `mapstructure:"dsn"` // mysql:// or postgres:// (目前仅支持 mysql 和 postgres)
	LogMode         bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

var C = Conf{
	DSN:             "postgres://postgres:postgres@localhost/test",
	LogMode:         false,
	MaxIdleConns:    0,
	MaxOpenConns:    0,
	ConnMaxLifetime: 0,
}

// BindDatabaseFlags 绑定 database viper.BingPFlag 的 key 增加绑定前缀
// @param keyPrefix key 前缀
// example:
// # config
// type Config struct {
//     Database Conf
// }
// BindDatabaseFlags(cmd, "database")
func BindDatabaseFlags(flagSet *flag.FlagSet, keyPrefix string) {
	flagSet.String("database_dsn", "", "database dsn eg: postgres://username:password@localhost/dbname")
	_ = viper.BindPFlag(keyPrefix+".DSN", flagSet.Lookup("database_dsn"))
	flagSet.Bool("database_log_mode", false, "gorm log mode")
	_ = viper.BindPFlag(keyPrefix+".LogMode", flagSet.Lookup("database_log_mode"))
	flagSet.Int("database_max_idle_conns", 0, "database max idle conns")
	_ = viper.BindPFlag(keyPrefix+".MaxIdleConns", flagSet.Lookup("database_max_idle_conns"))
	flagSet.Int("database_max_open_conns", 0, "database max open conns")
	_ = viper.BindPFlag(keyPrefix+".MaxOpenConns", flagSet.Lookup("database_max_open_conns"))
	flagSet.Int("database_conn_max_lifetime", 0, "database max open conns")
	_ = viper.BindPFlag(keyPrefix+".ConnMaxLifetime", flagSet.Lookup("database_conn_max_lifetime"))
}
