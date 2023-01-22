package conf

import env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type ExternalAppServerConf interface {
	GetUserServiceAddress() string
	GetKeyServiceAddress() string
	GetNoteServiceAddress() string
}

type ExternalAppServerConfImpl struct {
	userServiceAddress string
	keyServiceAddress  string
	noteServiceAddress string
}

func (e ExternalAppServerConfImpl) GetUserServiceAddress() string {
	return e.userServiceAddress
}

func (e ExternalAppServerConfImpl) GetKeyServiceAddress() string {
	return e.keyServiceAddress
}

func (e ExternalAppServerConfImpl) GetNoteServiceAddress() string {
	return e.noteServiceAddress
}

func NewExternalAppServerConfImpl() *ExternalAppServerConfImpl {
	return &ExternalAppServerConfImpl{
		userServiceAddress: env.GetEnvVar(env.EnvVarKeyAppserverUserServiceAddress),
		keyServiceAddress:  env.GetEnvVar(env.EnvVarKeyAppserverKeyServiceAddress),
		noteServiceAddress: env.GetEnvVar(env.EnvVarKeyAppserverNoteServiceAddress),
	}
}
