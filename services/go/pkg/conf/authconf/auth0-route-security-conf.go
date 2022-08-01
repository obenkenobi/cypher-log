package authconf

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Auth0RouteSecurityConf interface {
	GetIssuerUrl() *url.URL
	GetAudience() string
}

type Auth0RouteSecurityConfImpl struct {
	issuerUrl *url.URL
	audience  string
}

func (a Auth0RouteSecurityConfImpl) GetIssuerUrl() *url.URL {
	return a.issuerUrl
}

func (a Auth0RouteSecurityConfImpl) GetAudience() string {
	return a.audience
}

func NewAuth0RouteSecurityConf(envVarKeyAuth0IssuerUrl, envVarKeyAuth0Audience string) Auth0RouteSecurityConf {
	issuerUrlStr := environment.GetEnvVariable(envVarKeyAuth0IssuerUrl)
	issuerUrl, err := url.Parse(issuerUrlStr)
	if err != nil {
		log.Fatalf("Failed to parse issuer url %v", issuerUrlStr)
	}
	return &Auth0RouteSecurityConfImpl{
		issuerUrl: issuerUrl,
		audience:  environment.GetEnvVariable(envVarKeyAuth0Audience),
	}
}
