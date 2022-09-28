package authconf

import (
	"fmt"
	environment2 "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
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
}

type Auth0RouteSecurityConfImpl struct {
	issuerUrl               *url.URL
	apiAudience             string
	grpcAudience            string
	domain                  string
	clientCredentialsId     string
	clientCredentialsSecret string
}

func (a Auth0RouteSecurityConfImpl) GetGrpcAudience() string {
	return a.grpcAudience
}

func (a Auth0RouteSecurityConfImpl) GetIssuerUrl() *url.URL {
	return a.issuerUrl
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

func (a Auth0RouteSecurityConfImpl) GetApiAudience() string {
	return a.apiAudience
}

func NewAuth0SecurityConf() Auth0SecurityConf {
	issuerUrlStr := fmt.Sprintf("https://%v/", environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Domain))
	issuerUrl, err := url.Parse(issuerUrlStr)
	if err != nil {
		logger.Log.Fatalf("Failed to parse issuer url %v", issuerUrlStr)
	}
	return &Auth0RouteSecurityConfImpl{
		issuerUrl:               issuerUrl,
		apiAudience:             environment2.GetEnvVariable(environment2.EnvVarKeyAuth0ApiAudience),
		grpcAudience:            environment2.GetEnvVariable(environment2.EnvVarKeyAuth0GrpcAudience),
		domain:                  environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Domain),
		clientCredentialsId:     environment2.GetEnvVariable(environment2.EnvVarKeyAuth0ClientCredentialsId),
		clientCredentialsSecret: environment2.GetEnvVariable(environment2.EnvVarKeyAuth0ClientCredentialsSecret),
	}
}
