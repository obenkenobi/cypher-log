package authconf

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"net/url"
)

type Auth0SecurityConf interface {
	GetIssuerUrl() *url.URL
	GetApiAudience() string
	GetGrpcAudience() string
	GetDomain() string
	GetClientCredentialsId() string
	GetClientCredentialsSecret() string
	GetWebappClientId() string
	GetWebappClientSecret() string
	GetWebappCallbackUrl() string
}

type Auth0RouteSecurityConfImpl struct {
	issuerUrl               *url.URL
	apiAudience             string
	grpcAudience            string
	domain                  string
	clientCredentialsId     string
	clientCredentialsSecret string
	webappClientId          string
	webappClientSecret      string
	webappCallbackUrl       string
}

func (a Auth0RouteSecurityConfImpl) GetGrpcAudience() string {
	return a.grpcAudience
}

func (a Auth0RouteSecurityConfImpl) GetIssuerUrl() *url.URL {
	return a.issuerUrl
}

func (a Auth0RouteSecurityConfImpl) GetApiAudience() string {
	return a.apiAudience
}

func (a Auth0RouteSecurityConfImpl) GetDomain() string {
	return a.domain
}

func (a Auth0RouteSecurityConfImpl) GetClientCredentialsId() string {
	return a.clientCredentialsId
}

func (a Auth0RouteSecurityConfImpl) GetClientCredentialsSecret() string {
	return a.clientCredentialsSecret
}

func (a Auth0RouteSecurityConfImpl) GetWebappClientId() string {
	return a.webappClientId
}

func (a Auth0RouteSecurityConfImpl) GetWebappClientSecret() string {
	return a.webappClientSecret
}

func (a Auth0RouteSecurityConfImpl) GetWebappCallbackUrl() string {
	return a.webappCallbackUrl
}

func NewAuth0SecurityConfImpl() *Auth0RouteSecurityConfImpl {
	issuerUrlStr := fmt.Sprintf("https://%v/", environment.GetEnvVariable(environment.EnvVarKeyAuth0Domain))
	issuerUrl, err := url.Parse(issuerUrlStr)
	if err != nil {
		logger.Log.Fatalf("Failed to parse issuer url %v", issuerUrlStr)
	}
	return &Auth0RouteSecurityConfImpl{
		issuerUrl:               issuerUrl,
		apiAudience:             environment.GetEnvVariable(environment.EnvVarKeyAuth0ApiAudience),
		grpcAudience:            environment.GetEnvVariable(environment.EnvVarKeyAuth0GrpcAudience),
		domain:                  environment.GetEnvVariable(environment.EnvVarKeyAuth0Domain),
		clientCredentialsId:     environment.GetEnvVariable(environment.EnvVarKeyAuth0ClientCredentialsId),
		clientCredentialsSecret: environment.GetEnvVariable(environment.EnvVarKeyAuth0ClientCredentialsSecret),
		webappClientId:          environment.GetEnvVariable(environment.EnvVarKeyAuth0WebappClientId),
		webappClientSecret:      environment.GetEnvVariable(environment.EnvVarKeyAuth0WebappClientSecret),
		webappCallbackUrl:       environment.GetEnvVariable(environment.EnvVarKeyAuth0WebappCallbackUrl),
	}
}
