package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type TLSConf interface {
	ServerCertPath() string
	ServerKeyPath() string
	CACertPath() string
	WillLoadCACert() bool
}

type tlsConfImpl struct {
	serverCertPath string
	serverKeyPath  string
	caCertPath     string
	loadCACert     bool
}

func (t tlsConfImpl) ServerCertPath() string { return t.serverCertPath }

func (t tlsConfImpl) ServerKeyPath() string { return t.serverKeyPath }

func (t tlsConfImpl) CACertPath() string { return t.caCertPath }

func (t tlsConfImpl) WillLoadCACert() bool { return t.loadCACert }

func NewTlsConf() TLSConf {
	return &tlsConfImpl{
		serverCertPath: env.GetEnvVariable(env.EnvVarKeyServerCertPath),
		serverKeyPath:  env.GetEnvVariable(env.EnvVarKeyServerKeyPath),
		caCertPath:     env.GetEnvVariable(env.EnvVarKeyCACertPath),
		loadCACert:     env.GetEnvVarAsBoolOrDefault(env.EnvVarLoadCACert, false),
	}
}
