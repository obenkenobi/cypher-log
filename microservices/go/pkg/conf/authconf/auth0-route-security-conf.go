package authconf

import (
	"fmt"
	environment2 "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
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

func NewAuth0RouteSecurityConf() Auth0RouteSecurityConf {
	issuerUrlStr := fmt.Sprintf("https://%v/", environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Domain))
	issuerUrl, err := url.Parse(issuerUrlStr)
	if err != nil {
		log.Fatalf("Failed to parse issuer url %v", issuerUrlStr)
	}
	return &Auth0RouteSecurityConfImpl{
		issuerUrl: issuerUrl,
		audience:  environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Audience),
	}
}
