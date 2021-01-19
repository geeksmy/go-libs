//go:generate mockgen --source jwt.go --destination jwt.mock.go --package=jwt
// package jwt 签发和验签接口以及对应的 mock
package jwt

import (
	"crypto/rsa"
	"errors"
	"sync"

	jwtlib "github.com/dgrijalva/jwt-go"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// warn(joe@2019/11/19): 这里的默认配置被设计为只能用于资源接口请求认证, 用作其他用途可能带来未知的安全风险
type jwtOption struct {
	hmacSecret    []byte
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey

	tokenIssuer string
}

const (
	defaultTokenIssuer = "authority"
)

type SigningMethod string

const (
	SigningMethodRS256 SigningMethod = "RS256"
	SigningMethodRS512 SigningMethod = "RS512"
	SigningMethodHS512 SigningMethod = "HS512"
)

func (sm SigningMethod) getSigningMethod() jwtlib.SigningMethod {
	return jwtlib.GetSigningMethod(string(sm))
}

var (
	sharedOption *jwtOption

	ErrJwtSecretNotConfig = errors.New("jwt secret not config")
	ErrNotSupportedClaims = errors.New("only support userClaims")
)

// 使用默认配置初始化
func Init() error {
	return SetupWithConf(C)
}

// 使用指定配置初始化
func SetupWithConf(conf Conf) error {
	var err error
	sharedOption, err = newJwtOption(conf)

	return err
}

func newAccessList() ([]*AccessListEntry, error) {
	entry := NewAccessListEntry()
	entry.Allow()
	if err := entry.SetClaim("roles"); err != nil {
		return nil, ErrInvalidConfiguration.WithArgs("jwt", err)
	}

	for _, v := range []string{"anonymous", "guest"} {
		if err := entry.AddValue(v); err != nil {
			return nil, ErrInvalidConfiguration.WithArgs("jwt", err)
		}
	}
	accessList := make([]*AccessListEntry, 0, 1)
	accessList = append(accessList, entry)

	return accessList, nil
}

func newJwtOption(conf Conf) (*jwtOption, error) {
	option := &jwtOption{}

	if err := option.initWithConf(conf); err != nil {
		return nil, err
	}

	return option, nil
}

// initOptionWithConf 设置 jwt secret 和 keypair
// RSA 方式下签发方只需要设置 priKey 验签方只需要设置 Pubkey 即可
func (o *jwtOption) initWithConf(conf Conf) error {
	if err := conf.Validate(); err != nil {
		return err
	}

	if conf.Secret != "" {
		o.setHmacSecret(conf.Secret)
	}

	if conf.RSAPrivateKey != "" {
		if err := o.setRsaPrivateKey(conf.RSAPrivateKey); err != nil {
			return ErrInvaliRSAPrivateKey.WithMessage(err.Error())
		}
	}

	if conf.RSAPublicKey != "" {
		if err := o.setRSAPublicKey(conf.RSAPublicKey); err != nil {
			return ErrInvalidRSAPublicKey.WithMessage(err.Error())
		}
	}

	if conf.TokenIssuer != "" {
		o.tokenIssuer = conf.TokenIssuer
	} else {
		o.tokenIssuer = defaultTokenIssuer
	}

	return nil
}

func (o *jwtOption) setHmacSecret(secret string) {
	o.hmacSecret = []byte(secret)
}

func (o *jwtOption) setRsaPrivateKey(key string) error {
	priKey, err := ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return err
	}

	o.rsaPrivateKey = priKey
	return nil
}

func (o *jwtOption) setRSAPublicKey(key string) error {
	pubKey, err := ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return err
	}

	o.rsaPublicKey = pubKey
	return nil
}

// NewSigner 方便调用方可以 mock
// ctrl := gomock.NewCtrl(t)
// mocked := NewMockSigner(ctrl)
// jwt.NewSigner = func() Signer { return mocked }
// signer := jwt.NewSigner()
var NewSigner = NewSignerImpl

// NewValidator .
var NewValidator = NewValidatorImpl

// Signer jwt signer
type Signer interface {
	// Sign 签发 jwt
	Sign(claim jwtlib.Claims) (string, error)
}

// Validator jwt validator interface
type Validator interface {
	// 验证 jwt 是否合法
	Verify(tokenStr string, scheme *security.JWTScheme) (*UserClaims, error)
}

func NewSignerImpl() Signer {
	if sharedOption == nil {
		panic("jwt 未初始化")
	}
	return SignerImpl{
		option: sharedOption,
	}
}

func NewSignerImplWithConf(c Conf) (Signer, error) {
	option, err := newJwtOption(c)
	if err != nil {
		return nil, err
	}

	return SignerImpl{
		option: option,
	}, nil
}

// SignerImpl `Signer` 的默认实现
// warn(joe@2019/11/19): 默认实现的 Signer 设计用于登陆认证无法作其他用途, 用作其他用途可能带来未知的安全风险
type SignerImpl struct {
	option *jwtOption
}

func (s SignerImpl) Sign(claim jwtlib.Claims) (string, error) {
	userClaims, ok := claim.(UserClaims)
	if !ok {
		return "", ErrNotSupportedClaims
	}

	// 优先使用 RSA 签名
	if s.option.rsaPrivateKey != nil {
		return userClaims.GetToken(SigningMethodRS512, s.option.rsaPrivateKey)
	}

	if s.option.hmacSecret != nil {
		return userClaims.GetToken(SigningMethodHS512, s.option.hmacSecret)
	}

	return "", ErrJwtSecretNotConfig
}

func NewValidatorImpl() Validator {
	if sharedOption == nil {
		panic("jwt 未初始化")
	}
	impl := ValidatorImpl{
		option: sharedOption,
	}

	if err := impl.setup(); err != nil {
		panic("jwt 未初始化")
	}

	return impl
}

func NewValidatorImplWithConf(c Conf) (Validator, error) {
	option, err := newJwtOption(c)
	if err != nil {
		return nil, err
	}

	impl := ValidatorImpl{
		option: option,
	}

	if err := impl.setup(); err != nil {
		return nil, err
	}

	return impl, nil
}

// ValidatorImpl warn(joe@2019/11/19): 这个 validator 只能用来验证由 `SingerImpl` 签发的 jwt
type ValidatorImpl struct {
	option *jwtOption

	tokenValidator *TokenValidator
}

func (v *ValidatorImpl) setup() error {
	v.tokenValidator = NewTokenValidator()

	v.tokenValidator.TokenSecret = string(v.option.hmacSecret)
	v.tokenValidator.PrivateKey = v.option.rsaPrivateKey
	v.tokenValidator.PublicKey = v.option.rsaPublicKey
	v.tokenValidator.TokenIssuer = v.option.tokenIssuer

	v.tokenValidator.AccessList, _ = newAccessList()

	if err := v.tokenValidator.ConfigureTokenBackends(); err != nil {
		return ErrInvalidBackendConfiguration.WithArgs("jwt", err)
	}

	return nil
}

func (v ValidatorImpl) Verify(token string, scheme *security.JWTScheme) (*UserClaims, error) {
	userClaims, valid, err := v.tokenValidator.ValidateToken(token)
	if err != nil {
		// nolint(errorlint): fixme
		if e, ok := err.(*jwtlib.ValidationError); ok {
			return nil, UnauthorizedErr(translateJwtValidationError(e))
		}
		return nil, err
	}

	if !valid {
		return nil, UnauthorizedErr("令牌错误")
	}

	return userClaims, nil
}

func UnauthorizedErr(format string, args ...interface{}) error {
	return goa.TemporaryError("unauthorized", format, args...)
}

var (
	validator     Validator
	signer        Signer
	validatorOnce sync.Once
	signerOnce    sync.Once
)

// ValidateToken returns a UserClaims for JWT token
func ValidateToken(token string, scheme *security.JWTScheme) (*UserClaims, error) {
	validatorOnce.Do(func() {
		validator = NewValidator()
	})

	return validator.Verify(token, scheme)
}

// GetToken returns a signed JWT token
func GetToken(claims jwtlib.Claims) (string, error) {
	signerOnce.Do(func() {
		signer = NewSigner()
	})

	return signer.Sign(claims)
}
