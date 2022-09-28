package externalservices

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
)

type HttpClientProvider interface {
	Client() *retryablehttp.Client
}

type HttpClientProviderImpl struct {
	httpClientConf conf.HttpClientConf
}

func (h HttpClientProviderImpl) Client() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = h.httpClientConf.GetRetryCount()
	return retryClient
}

func NewHTTPClientProvider(httpClientConf conf.HttpClientConf) HttpClientProvider {
	return &HttpClientProviderImpl{httpClientConf: httpClientConf}
}
