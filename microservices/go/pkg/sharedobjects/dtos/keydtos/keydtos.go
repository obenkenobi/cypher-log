package keydtos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"
)

type PasscodeCreateDto struct {
	Passcode string `json:"passcode" binding:"exists,alphanumunicode,min=4,max=20"`
}

type PasscodeDto struct {
	Passcode string `json:"passcode" binding:"required"`
}

type UserKeySessionTokenDto struct {
	ProxyKid string `json:"proxyKid" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type UserKeyDto struct {
	KeyBase64 string `json:"key"`
}

func NewUserKeyDto(key []byte) UserKeyDto {
	return UserKeyDto{
		KeyBase64: encodingutils.EncodeBase64String(key),
	}
}

func (u UserKeyDto) GetKey() ([]byte, error) {
	return encodingutils.DecodeBase64String(u.KeyBase64)
}
