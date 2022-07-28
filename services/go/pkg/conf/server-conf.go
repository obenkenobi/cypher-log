package conf

import "github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"

type ServerConf interface {
	GetPort() string
}

type serverConfImpl struct {
	port string
}

func (s serverConfImpl) GetPort() string {
	return s.port
}

func NewServerConf(envVarKeyPort string) ServerConf {
	return &serverConfImpl{
		port: environment.GetEnvVariableOrDefault(envVarKeyPort, "8080"),
	}
}
