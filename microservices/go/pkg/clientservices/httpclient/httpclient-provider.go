package httpclient

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
)

type ClientProvider interface {
	Client() *retryablehttp.Client
}

type ClientProviderImpl struct {
	httpClientConf conf.HttpClientConf
}

func (h ClientProviderImpl) Client() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = h.httpClientConf.GetRetryCount()
	return retryClient
}

func NewClientProvider(httpClientConf conf.HttpClientConf) *ClientProviderImpl {
	return &ClientProviderImpl{httpClientConf: httpClientConf}
}
