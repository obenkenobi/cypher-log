package userdtos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
)

type UserIdentityDto struct {
	UserReadDto
	embeddeduser.BaseUserAuthId
	Authorities []string `json:"authorities"`
}

type UserReadDto struct {
	embeddeduser.BaseUserPublicDto
	Exists bool `json:"exists"`
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

type UserChangeEventResponseDto struct {
	Discarded bool
}

type UserSaveDto struct {
	embeddeduser.BaseUserCommon
}
