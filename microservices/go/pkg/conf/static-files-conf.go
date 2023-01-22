package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type StaticFilesConf interface {
	GetStaticFilesPath() string
}

type StaticFilesConfImpl struct {
	staticFilesPath string
}

func (s StaticFilesConfImpl) GetStaticFilesPath() string {
	return s.staticFilesPath
}

func NewStaticFilesConfImpl() *StaticFilesConfImpl {
	return &StaticFilesConfImpl{
		staticFilesPath: env.GetEnvVar(env.EnvVarStaticFilesPath),
	}
}
