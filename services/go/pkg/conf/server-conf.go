package conf

type ServerConf interface {
	GetPort() string
}

type serverConfImpl struct {
	port string
}

func (s serverConfImpl) GetPort() string {
	return s.port
}

func NewServerConf(envVarReader EnvVarReader, envVarKeyPort string) ServerConf {
	return &serverConfImpl{
		port: envVarReader.GetEnvVariableOrDefault(envVarKeyPort, "8080"),
	}
}
