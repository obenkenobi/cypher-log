package conf

import "time"

type KeyConf interface {
	GetTokenSessionDuration() time.Duration
	GetSecretDuration() time.Duration
	GetKeyRefreshInterval() time.Duration
}

type keyConfImpl struct {
	tokenSessionDuration time.Duration
	secretDuration       time.Duration
	keyRefreshInterval   time.Duration
}

func NewKeyConfImpl() *keyConfImpl {
	tokenSessionDuration := 30 * time.Minute
	secretDuration := 3 * tokenSessionDuration
	keyRefreshInterval := 2 * tokenSessionDuration
	return &keyConfImpl{
		tokenSessionDuration: tokenSessionDuration,
		secretDuration:       secretDuration,
		keyRefreshInterval:   keyRefreshInterval,
	}
}
