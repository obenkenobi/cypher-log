package securityservices

import (
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"net/url"
	"time"
)

type BaseJwtValidateService interface {
	GetJwtValidator() *validator.Validator
}

type BaseOath2ValidateServiceImpl struct {
	provider     *jwks.CachingProvider
	jwtValidator *validator.Validator
}

func (o BaseOath2ValidateServiceImpl) GetJwtValidator() *validator.Validator {
	return o.jwtValidator
}

func createOath2ValidateService(issuerURL *url.URL, audience string) *BaseOath2ValidateServiceImpl {
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(func() validator.CustomClaims { return &security.Auth0CustomClaims{} }),
		validator.WithAllowedClockSkew(time.Minute),
	)
	return &BaseOath2ValidateServiceImpl{
		provider:     provider,
		jwtValidator: jwtValidator,
	}
}

type JwtValidateWebAppService interface {
	BaseJwtValidateService
}

type JwtValidateGrpcService interface {
	BaseJwtValidateService
}

type JwtValidateWebAppServiceImpl struct {
	BaseJwtValidateService
}

type JwtValidateGrpcServiceImpl struct {
	BaseJwtValidateService
}

func NewJwtValidateWebAppServiceImpl(
	auth0RouteSecurityConf authconf.Auth0SecurityConf,
) *JwtValidateWebAppServiceImpl {
	baseService := createOath2ValidateService(
		auth0RouteSecurityConf.GetIssuerUrl(),
		auth0RouteSecurityConf.GetApiAudience(),
	)
	return &JwtValidateWebAppServiceImpl{BaseJwtValidateService: baseService}
}

func NewJwtValidateGrpcServiceImpl(
	auth0RouteSecurityConf authconf.Auth0SecurityConf,
) *JwtValidateGrpcServiceImpl {
	baseService := createOath2ValidateService(
		auth0RouteSecurityConf.GetIssuerUrl(),
		auth0RouteSecurityConf.GetGrpcAudience(),
	)
	return &JwtValidateGrpcServiceImpl{BaseJwtValidateService: baseService}
}
