package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type TLSConf interface {
	ServerCertPath() string
	ServerKeyPath() string
	CACertPath() string
	WillLoadCACert() bool
}

type TlsConfImpl struct {
	serverCertPath string
	serverKeyPath  string
	caCertPath     string
	loadCACert     bool
}

func (t TlsConfImpl) ServerCertPath() string { return t.serverCertPath }

func (t TlsConfImpl) ServerKeyPath() string { return t.serverKeyPath }

func (t TlsConfImpl) CACertPath() string { return t.caCertPath }

func (t TlsConfImpl) WillLoadCACert() bool { return t.loadCACert }

func NewTlsConfImpl() *TlsConfImpl {
	return &TlsConfImpl{
		serverCertPath: env.GetEnvVariable(env.EnvVarKeyServerCertPath),
		serverKeyPath:  env.GetEnvVariable(env.EnvVarKeyServerKeyPath),
		caCertPath:     env.GetEnvVariable(env.EnvVarKeyCACertPath),
		loadCACert:     env.GetEnvVarAsBoolOrDefault(env.EnvVarLoadCACert, false),
	}
}
