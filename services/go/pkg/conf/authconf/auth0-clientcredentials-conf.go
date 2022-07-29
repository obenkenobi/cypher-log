package authconf

import "github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"

type Auth0ClientCredentialsConf interface {
	GetDomain() string
	GetClientId() string
	GetClientSecret() string
	GetAudience() string
}

type Auth0ClientCredentialsConfImpl struct {
	domain       string
	clientId     string
	clientSecret string
	audience     string
}

func (a Auth0ClientCredentialsConfImpl) GetDomain() string {
	return a.domain
}

func (a Auth0ClientCredentialsConfImpl) GetClientId() string {
	return a.clientId
}

func (a Auth0ClientCredentialsConfImpl) GetClientSecret() string {
	return a.clientSecret
}

func (a Auth0ClientCredentialsConfImpl) GetAudience() string {
	return a.audience
}

func NewAuth0ClientCredentialsConf(
	envVarKeyAuth0Domain string,
	envAuth0ClientId string,
	envAuth0ClientSecret string,
	envAuth0Audience string,
) Auth0ClientCredentialsConf {
	return &Auth0ClientCredentialsConfImpl{
		domain:       environment.GetEnvVariable(envVarKeyAuth0Domain),
		clientId:     environment.GetEnvVariable(envAuth0ClientId),
		clientSecret: environment.GetEnvVariable(envAuth0ClientSecret),
		audience:     environment.GetEnvVariable(envAuth0Audience),
	}
}
