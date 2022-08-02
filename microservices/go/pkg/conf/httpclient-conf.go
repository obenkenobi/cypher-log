package conf

type HttpClientConf interface {
	GetRetryCount() int
	GetJSONContentType() string
}

type HttpClientConfImpl struct {
}

func (h HttpClientConfImpl) GetRetryCount() int { return 10 }

func (h HttpClientConfImpl) GetJSONContentType() string { return "application/json" }

func NewHttpClientConfImpl() *HttpClientConfImpl {
	return &HttpClientConfImpl{}
}
