package conf

import "time"

type KeyConf interface {
	GetTokenSessionDuration() time.Duration
	GetSecretDuration() time.Duration
	GetKeyRefreshInterval() time.Duration
	GetPrimaryAppSecretDuration() time.Duration
}

type KeyConfImpl struct {
	tokenSessionDuration     time.Duration
	secretDuration           time.Duration
	keyRefreshInterval       time.Duration
	primaryAppSecretDuration time.Duration
}

func (k KeyConfImpl) GetTokenSessionDuration() time.Duration {
	return k.tokenSessionDuration
}

func (k KeyConfImpl) GetSecretDuration() time.Duration {
	return k.secretDuration
}

func (k KeyConfImpl) GetKeyRefreshInterval() time.Duration {
	return k.keyRefreshInterval
}

func (k KeyConfImpl) GetPrimaryAppSecretDuration() time.Duration {
	return k.primaryAppSecretDuration
}

func NewKeyConfImpl() *KeyConfImpl {
	tokenSessionDuration := 30 * time.Minute
	secretDuration := 6 * tokenSessionDuration
	keyRefreshInterval := 3 * tokenSessionDuration
	primaryAppSecretDuration := 4 * tokenSessionDuration
	return &KeyConfImpl{
		tokenSessionDuration:     tokenSessionDuration,
		secretDuration:           secretDuration,
		keyRefreshInterval:       keyRefreshInterval,
		primaryAppSecretDuration: primaryAppSecretDuration,
	}
}
