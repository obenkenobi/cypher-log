package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type KafkaConf interface {
	GetBootstrapServers() []string
	GetUsername() string
	GetPassword() string
}

type KafkaConfImpl struct {
	Servers  []string
	Username string
	Password string
}

func (k KafkaConfImpl) GetBootstrapServers() []string {
	return k.Servers
}

func (k KafkaConfImpl) GetUsername() string {
	return k.Username
}

func (k KafkaConfImpl) GetPassword() string {
	return k.Password
}

func NewKafkaConfImpl() *KafkaConfImpl {
	return &KafkaConfImpl{
		Servers:  env.GetEnvVariableAsListSplitByComma(env.EnvVarKafkaBootstrapServers),
		Username: env.GetEnvVar(env.EnvVarKafkaUsername),
		Password: env.GetEnvVar(env.EnvVarKafkaPassword),
	}
}
