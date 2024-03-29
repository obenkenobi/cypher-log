package userdtos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
)

type UserIdentityDto struct {
	UserReadDto
	embeddeduser.BaseUserAuthId
	Authorities []string `json:"authorities"`
}

type UserReadDto struct {
	embeddeduser.BaseUserPublicDto
	commondtos.ExistsDto
}

type UserChangeAction int64

const (
	UserSave UserChangeAction = iota
	UserDelete
)

type UserChangeEventDto struct {
	embeddeduser.BaseUserPublicDto
	embeddeduser.BaseUserAuthId
	Action UserChangeAction `json:"action"`
}

func (u UserChangeEventDto) MessageKey() ([]byte, error) {
	return []byte(u.Id), nil
}

type UserChangeEventResponseDto struct {
	Discarded bool
}

type UserSaveDto struct {
	embeddeduser.BaseUserCommon
}
