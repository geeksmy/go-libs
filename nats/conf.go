package nats

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
# 默认配置
NAT:
  # server 地址
  Host: "nats://127.0.0.1:4222"
  Username: ""
  Password: ""
  # 连接名称
  ClientName: "nat-client"
*/
type Conf struct {
	Host       string // server地址，逗号分隔
	Username   string
	Password   string
	ClientName string
}

var C = Conf{
	Host:       "nats://127.0.0.1:4222",
	Username:   "",
	Password:   "",
	ClientName: "nat-client",
}

func BindNatFlags(flagSet *flag.FlagSet, keyPrefix string) {
	flagSet.String("nat_host", "", "nat server host, comma-separated")
	_ = viper.BindPFlag(keyPrefix+".Host", flagSet.Lookup("nat_host"))
	flagSet.String("nat_username", "", "username for auth")
	_ = viper.BindPFlag(keyPrefix+".Username", flagSet.Lookup("nat_username"))
	flagSet.String("nat_password", "", "username for auth")
	_ = viper.BindPFlag(keyPrefix+".Password", flagSet.Lookup("nat_password"))
	flagSet.String("nat_client_name", "", "client name")
	_ = viper.BindPFlag(keyPrefix+".ClientName", flagSet.Lookup("nat_client_name"))
}
