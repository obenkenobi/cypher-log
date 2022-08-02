package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/clientcredentialsdtos"
)

type SystemAccessTokenClient interface {
	getAccessToken() (string, error)
}

type Auth0SystemAccessTokenClient struct {
	httpClientConf        conf.HttpClientConf
	clientCredentialsConf authconf.Auth0ClientCredentialsConf
	clientProvider        ClientProvider
}

func (a Auth0SystemAccessTokenClient) getAccessToken() (string, error) {
	token := ""
	url := fmt.Sprintf("https://%v/oauth/token", a.clientCredentialsConf.GetDomain())
	client := a.clientProvider.Client()
	payload := clientcredentialsdtos.ClientCredentialsRequestDto{
		ClientId:     a.clientCredentialsConf.GetClientId(),
		ClientSecret: a.clientCredentialsConf.GetClientSecret(),
		Audience:     a.clientCredentialsConf.GetAudience(),
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
