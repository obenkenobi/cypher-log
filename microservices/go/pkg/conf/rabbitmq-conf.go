package conf

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type RabbitMQConf interface {
	GetURI() string
}

type RabbitMQConfImpl struct {
	uri string
}

func (r RabbitMQConfImpl) GetURI() string { return r.uri }

func NewRabbitMQConfImpl() *RabbitMQConfImpl {
	return &RabbitMQConfImpl{uri: environment.GetEnvVar(environment.EnvVarRabbitMQUri)}
}
