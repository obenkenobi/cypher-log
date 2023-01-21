package security

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New(auth0SecurityConf authconf.Auth0SecurityConf) *Authenticator {
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

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
