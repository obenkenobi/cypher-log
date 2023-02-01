package keydtos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"
)

type PasscodeCreateDto struct {
	Passcode string `json:"passcode" binding:"required,alphanumunicode,min=4,max=20"`
}

type PasscodeDto struct {
	Passcode string `json:"passcode" binding:"required"`
}

type UserKeyDto struct {
	KeyBase64  string `json:"keyBase64"`
	KeyVersion int64  `json:"keyVersion"`
}

func NewUserKeyDto(keyBytes []byte, keyVersion int64) UserKeyDto {
	return UserKeyDto{
		KeyBase64:  encodingutils.EncodeBase64String(keyBytes),
		KeyVersion: keyVersion,
	}
}

func (u UserKeyDto) GetKey() ([]byte, error) {
	return encodingutils.DecodeBase64String(u.KeyBase64)
}
