// 提供 goa Server 的配置方式
package goalibs

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
Serve:
    Host: localhost
	Domain: 127.0.0.1
	HTTPPort: 8080
	GRPCPort: 8081
	Secure: false
	Debug: false
*/
// goa example server 配置
type ServeConf struct {
	Host     string
	Domain   string
	HTTPPort int
	GRPCPort int
	Secure   bool
	Debug    bool
}

// 为 serve cmd 和 Serve 绑定 pflag
func BindPflag(flagSet *pflag.FlagSet, keyPrefix string) {
	flagSet.String("host", "", "Server host")
	_ = viper.BindPFlag(keyPrefix+".host", flagSet.Lookup("host"))

	flagSet.String("domain", "", "Host domain name (overrides host domain specified in service design)")
	_ = viper.BindPFlag(keyPrefix+".domain", flagSet.Lookup("domain"))

	flagSet.Int("http-port", 8080, "HTTP port (overrides host HTTP port specified in service design)")
	_ = viper.BindPFlag(keyPrefix+".HTTPPort", flagSet.Lookup("http-port"))

	flagSet.Int("grpc-port", 8082, "gRPC port (overrides host gRPC port specified in service design)")
	_ = viper.BindPFlag(keyPrefix+".GRPCPort", flagSet.Lookup("grpc-port"))

	flagSet.Bool("secure", false, "Use secure scheme (https or grpcs)")
	_ = viper.BindPFlag(keyPrefix+".secure", flagSet.Lookup("secure"))

	flagSet.Bool("debug", false, "Log request and response bodies")
	_ = viper.BindPFlag(keyPrefix+".debug", flagSet.Lookup("debug"))
}
