package externalservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/clientcredentialsdtos"
	"golang.org/x/oauth2"
)

type SysAccessTokenClient interface {
	GetApiAccessToken() (oauth2.Token, error)
	GetGRPCAccessToken() (oauth2.Token, error)
}

type Auth0SysAccessTokenClient struct {
	httpClientConf    conf.HttpClientConf
	auth0SecurityConf authconf.Auth0SecurityConf
	clientProvider    HttpClientProvider
}

func (a Auth0SysAccessTokenClient) GetApiAccessToken() (oauth2.Token, error) {
	return a.getAccessToken(a.auth0SecurityConf.GetApiAudience())
}

func (a Auth0SysAccessTokenClient) GetGRPCAccessToken() (oauth2.Token, error) {
	return a.getAccessToken(a.auth0SecurityConf.GetGrpcAudience())
}

func (a Auth0SysAccessTokenClient) getAccessToken(audience string) (oauth2.Token, error) {
	token := oauth2.Token{}
	url := fmt.Sprintf("https://%v/oauth/token", a.auth0SecurityConf.GetDomain())
	client := a.clientProvider.Client()
	payload := clientcredentialsdtos.ClientCredentialsRequestDto{
		ClientId:     a.auth0SecurityConf.GetClientCredentialsId(),
		ClientSecret: a.auth0SecurityConf.GetClientCredentialsSecret(),
		Audience:     audience,
		GrantType:    "client_credentials",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return token, err
	}
	resp, err := client.Post(url, a.httpClientConf.GetJSONContentType(), bytes.NewReader(jsonPayload))
	if err != nil {
		return token, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return token, err
	}
	return token, err
}

func NewAuth0SysAccessTokenClient(
	httpClientConf conf.HttpClientConf,
	auth0SecurityConf authconf.Auth0SecurityConf,
	clientProvider HttpClientProvider,
) SysAccessTokenClient {
	return &Auth0SysAccessTokenClient{
		httpClientConf:    httpClientConf,
		auth0SecurityConf: auth0SecurityConf,
		clientProvider:    clientProvider,
	}
}
