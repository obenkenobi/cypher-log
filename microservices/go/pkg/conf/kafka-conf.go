package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type KafkaConf interface {
	GetServers() []string
	GetUsername() string
	GetPassword() string
}

type KafkaConfImpl struct {
	Servers  []string
	Username string
	Password string
}

func (k KafkaConfImpl) GetServers() []string {
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
		Servers:  env.GetEnvVariableAsListSplitByComma(env.EnvVarKafkaServers),
		Username: env.GetEnvVar(env.EnvVarKafkaUsername),
		Password: env.GetEnvVar(env.EnvVarKafkaPassword),
	}
}
