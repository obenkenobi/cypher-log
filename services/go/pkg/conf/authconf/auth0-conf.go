package authconf

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Auth0Conf interface {
	GetIssuerUrl() *url.URL
	GetAudience() string
}

type Auth0ConfImpl struct {
	issuerUrl *url.URL
	audience  string
}

func (a Auth0ConfImpl) GetIssuerUrl() *url.URL {
	return a.issuerUrl
}

func (a Auth0ConfImpl) GetAudience() string {
	return a.audience
}

func NewAuth0Conf(envVarKeyAuth0IssuerUrl string, envVarKeyAuth0Audience string) *Auth0ConfImpl {
	issuerUrlStr := environment.GetEnvVariable(envVarKeyAuth0IssuerUrl)
	issuerUrl, err := url.Parse(issuerUrlStr)
	if err != nil {
		log.Fatalf("Failed to parse issuer url %v", issuerUrlStr)
	}
	return &Auth0ConfImpl{
		issuerUrl: issuerUrl,
		audience:  environment.GetEnvVariable(envVarKeyAuth0Audience),
	}
}
