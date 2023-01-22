package conf

import (
	env "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/randutils"
)

type SessionConf interface {
	GetSessionStoreSecret() string
}

type SessionConfImpl struct {
	sessionStoreSecret string
}

func (s SessionConfImpl) GetSessionStoreSecret() string {
	return s.sessionStoreSecret
}

func NewSessionConfImpl() *SessionConfImpl {
	defaultSessionStoreSecret, err := randutils.GenerateRandom32Bytes()
	if err != nil {
		logger.Log.WithError(err).Fatal()
	}
	return &SessionConfImpl{
		sessionStoreSecret: env.GetEnvVarOrDefault(env.EnvVarSessionStoreSecret, defaultSessionStoreSecret),
	}
}
