package securityservices

import (
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"net/url"
	"time"
)

type ExternalOath2ValidateService interface {
	GetJwtValidator() *validator.Validator
}

type ExternalOath2ValidateServiceImpl struct {
	provider     *jwks.CachingProvider
	jwtValidator *validator.Validator
}

func (o ExternalOath2ValidateServiceImpl) GetJwtValidator() *validator.Validator {
	return o.jwtValidator
}

func NewAPIAuth0JwtValidateService(auth0RouteSecurityConf authconf.Auth0SecurityConf) ExternalOath2ValidateService {
	return createOath2ValidateService(auth0RouteSecurityConf.GetIssuerUrl(), auth0RouteSecurityConf.GetApiAudience())
}

func NewGrpcAuth0JwtValidateService(auth0RouteSecurityConf authconf.Auth0SecurityConf) ExternalOath2ValidateService {
	return createOath2ValidateService(auth0RouteSecurityConf.GetIssuerUrl(), auth0RouteSecurityConf.GetGrpcAudience())
}

func createOath2ValidateService(issuerURL *url.URL, audience string) ExternalOath2ValidateService {
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(func() validator.CustomClaims { return &security.Auth0CustomClaims{} }),
		validator.WithAllowedClockSkew(time.Minute),
	)
	return &ExternalOath2ValidateServiceImpl{
		provider:     provider,
		jwtValidator: jwtValidator,
	}
}
