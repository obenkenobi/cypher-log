package conf

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type RabbitMQConf interface {
	GetURI() string
}

type rabbitMQConfImpl struct {
	uri string
}

func (r rabbitMQConfImpl) GetURI() string { return r.uri }

func NewRabbitMQConfImpl() *rabbitMQConfImpl {
	return &rabbitMQConfImpl{uri: environment.GetEnvVariable(environment.EnvVarRabbitMQUri)}
}
