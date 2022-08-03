package authconf

import (
	environment2 "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
)

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

func NewAuth0ClientCredentialsConf() Auth0ClientCredentialsConf {
	return &Auth0ClientCredentialsConfImpl{
		domain:       environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Domain),
		clientId:     environment2.GetEnvVariable(environment2.EnvVarKeyAuth0ClientId),
		clientSecret: environment2.GetEnvVariable(environment2.EnvVarKeyAuth0ClientSecret),
		audience:     environment2.GetEnvVariable(environment2.EnvVarKeyAuth0Audience),
	}
}
