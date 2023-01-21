package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"golang.org/x/oauth2"
)

type AuthenticatorService interface {
	GetOidcProvider() *oidc.Provider
	GetOath2Config() *oauth2.Config
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}

type AuthenticatorServiceImpl struct {
	provider *oidc.Provider
	config   *oauth2.Config
}

func (a AuthenticatorServiceImpl) GetOidcProvider() *oidc.Provider {
	return a.provider
}

func (a AuthenticatorServiceImpl) GetOath2Config() *oauth2.Config {
	return a.config
}

func (a AuthenticatorServiceImpl) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.GetOath2Config().ClientID,
	}

	return a.GetOidcProvider().Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func NewAuthenticatorServiceImpl(auth0SecurityConf authconf.Auth0SecurityConf) *AuthenticatorServiceImpl {
	provider, err := oidc.NewProvider(
		context.Background(),
		fmt.Sprintf("https://%v/", auth0SecurityConf.GetDomain()),
	)
	if err != nil {
		logger.Log.WithError(err).Fatal()
	}

	conf := oauth2.Config{
		ClientID:     auth0SecurityConf.GetWebappClientId(),
		ClientSecret: auth0SecurityConf.GetWebappClientSecret(),
		RedirectURL:  auth0SecurityConf.GetWebappCallbackUrl(),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &AuthenticatorServiceImpl{
		provider: provider,
		config:   &conf,
	}
}
