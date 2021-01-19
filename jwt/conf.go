package jwt

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Conf jwt 的配置项
// JWT:
//     Secret:
//     RSAPrivateKey: |
//         this is private key
//         multiline str
//	   RSAPublicKey: |
//         this is pub key
//         multiline str
type Conf struct {
	// 如果使用 HMAC 需要配置
	Secret string
	// RSA 的配置项目, 如果使用 rsa 签名/验签 需要配置
	RSAPrivateKey string
	RSAPublicKey  string
	TokenIssuer   string
}

const (
	ConfErrCodeNoRequiredKey = iota + 1
	ConfErrCodeRSAPrivateKey
	ConfErrCodeRSAPublic
)

type ConfErr struct {
	code    int
	message string
}

func (e ConfErr) Code() int {
	return e.code
}

// WithMessage clone ConfErr and use `msg` as `ConfErr.message`
func (e ConfErr) WithMessage(msg string) ConfErr {
	return ConfErr{
		code:    e.code,
		message: fmt.Sprintf("%s: %q", e.message, msg),
	}
}

func (e ConfErr) Error() string {
	return e.message
}

var (
	ErrNoRequiredSecret    = ConfErr{code: ConfErrCodeNoRequiredKey, message: "没配置 jwt 密钥"}
	ErrInvaliRSAPrivateKey = ConfErr{code: ConfErrCodeRSAPrivateKey, message: "jwt 私钥配置错误"}
	ErrInvalidRSAPublicKey = ConfErr{code: ConfErrCodeRSAPublic, message: "jwt 公钥配置错误"}
)

func (c Conf) Validate() error {
	if c.Secret == "" && !(c.RSAPrivateKey != "" || c.RSAPublicKey != "") {
		return ErrNoRequiredSecret
	}

	// if c.RSAPrivateKey != "" {
	// 	_, err := ParseRSAPrivateKeyFromPEM(c.RSAPrivateKey)
	// 	if err != nil {
	// 		return ErrInvaliRSAPrivateKey.WithMessage(err.Error())
	// 	}
	// }

	// if c.RSAPublicKey != "" {
	// 	_, err := ParseRSAPublicKeyFromPEM(c.RSAPublicKey)
	// 	if err != nil {
	// 		return ErrInvalidRSAPublicKey.WithMessage(err.Error())
	// 	}
	// }

	return nil
}

// C 默认配置项目, 调用方可以直接引用
// warn(joe@2019/11/19): 默认的配置被设计为只能用于资源请求接口授权验证, 用作其他用途可能带来未知的安全风险
// var Default = Config{Jwt: &jwt.C}
// err := jwt.Init()
var C = Conf{}

// BindPflag 为 jwt.Conf 提供 pflag 注册接口
func BindPflag(flagSet *pflag.FlagSet, keyPrefix string) {
	flagSet.String("jwt_secret", "", "jwt hmac secret")
	_ = viper.BindPFlag(keyPrefix+".Secret", flagSet.Lookup("jwt_secret"))

	flagSet.String("jwt_rsa_private_key", "", "jwt rsa private key")
	_ = viper.BindPFlag(keyPrefix+".RSAPrivateKey", flagSet.Lookup("jwt_rsa_private_key"))

	flagSet.String("jwt_rsa_public_key", "", "jwt rsa public key")
	_ = viper.BindPFlag(keyPrefix+".RSAPublicKey", flagSet.Lookup("jwt_rsa_public_key"))
}
