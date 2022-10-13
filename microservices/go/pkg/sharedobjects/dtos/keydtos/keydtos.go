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

type UserKeySessionDto struct {
	ProxyKid      string `json:"proxyKid" binding:"required"`
	Token         string `json:"token" binding:"required"`
	UserId        string `json:"userId" binding:"required"`
	KeyVersion    int64  `json:"keyVersion" binding:"required"`
	StartTime     int64  `json:"startTime"`     // In unix timestamp in milliseconds
	DurationMilli int64  `json:"durationMilli"` // In milliseconds
}

type UserKeyDto struct {
	KeyBase64 string `json:"key"`
}

func NewUserKeyDto(keyBytes []byte) UserKeyDto {
	return UserKeyDto{
		KeyBase64: encodingutils.EncodeBase64String(keyBytes),
	}
}

func (u UserKeyDto) GetKey() ([]byte, error) {
	return encodingutils.DecodeBase64String(u.KeyBase64)
}
