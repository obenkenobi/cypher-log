package models

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"

type AccessTokenHolder struct {
	TokenCipher []byte `json:"tokenCipher"`
}

func (a *AccessTokenHolder) SetAccessToken(accessToken string, key []byte) (err error) {
	a.TokenCipher, err = cipherutils.EncryptAES(key, []byte(accessToken))
	return err
}

func (a *AccessTokenHolder) GetAccessToken(key []byte) (string, error) {
	tokenBytes, err := cipherutils.DecryptAES(key, a.TokenCipher)
	return string(tokenBytes), err
}
