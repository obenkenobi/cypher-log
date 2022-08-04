package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/clientcredentialsdtos"
)

type SystemAccessTokenClient interface {
	GetApiAccessToken() (string, error)
	GetGRPCAccessToken() (string, error)
}

type Auth0SystemAccessTokenClient struct {
	httpClientConf    conf.HttpClientConf
	auth0SecurityConf authconf.Auth0SecurityConf
	clientProvider    ClientProvider
}

func (a Auth0SystemAccessTokenClient) GetApiAccessToken() (string, error) {
	return a.getAccessToken(a.auth0SecurityConf.GetApiAudience())
}

func (a Auth0SystemAccessTokenClient) GetGRPCAccessToken() (string, error) {
	return a.getAccessToken(a.auth0SecurityConf.GetGrpcAudience())
}

func (a Auth0SystemAccessTokenClient) getAccessToken(audience string) (string, error) {
	token := ""
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
	clientCredentialsResponse := clientcredentialsdtos.ClientCredentialsResultDto{}
	if err := json.NewDecoder(resp.Body).Decode(&clientCredentialsResponse); err != nil {
		return token, err
	}
	return clientCredentialsResponse.AccessToken, err
}

func NewAuth0SystemAccessTokenClient(
	httpClientConf conf.HttpClientConf,
	auth0SecurityConf authconf.Auth0SecurityConf,
	clientProvider ClientProvider,
) *Auth0SystemAccessTokenClient {
	return &Auth0SystemAccessTokenClient{
		httpClientConf:    httpClientConf,
		auth0SecurityConf: auth0SecurityConf,
		clientProvider:    clientProvider,
	}
}
